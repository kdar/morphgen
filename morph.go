package main

import (
  // "bitbucket.org/kardianos/osext"
  "errors"
  "flag"
  "fmt"
  "github.com/gosexy/to"
  "io"
  "os"
  // "path/filepath"
  "sort"
  "strings"
)

var (
  APPDIR string

  noTmog      = flag.Bool("notmog", false, "Toggle off grabbing transmogged items from armory")
  versionFlag = flag.Bool("version", false, "Display version")

  // a map from names to slots. this also contains
  // what we're able to tmorph into
  nameSlotMap = map[string]int{
    "head":     1,
    "shoulder": 3,
    "shirt":    4,
    "chest":    5,
    "waist":    6,
    "legs":     7,
    "feet":     8,
    "wrist":    9,
    "hands":    10,
    "back":     15,
    "mainHand": 16,
    "offHand":  17,
    "tabard":   19,
  }

  // mappings of items type ids to wow slot ids
  slotMap = map[int]int{
    // main hand slot
    13: 16, // one hand
    15: 16, // ranged
    17: 16, // two-hand
    21: 16, // main hand
    25: 16, // thrown

    // off hand slot
    14: 17, // shield
    22: 17, // off-hand
    23: 17, // held in off-hand
  }
)

type TMorphItem struct {
  Type string
  Args []int
}

func (t TMorphItem) String() string {
  args := ""
  for _, i := range t.Args {
    args = args + to.String(i) + " "
  }

  return fmt.Sprintf(".%s %s", t.Type, strings.Trim(args, " "))
}

type TMorphItems []*TMorphItem

func (t TMorphItems) Len() int      { return len(t) }
func (t TMorphItems) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t TMorphItems) Less(i, j int) bool {
  if t[i].Type != t[j].Type {
    // give priority to sorting items first.
    // just looks better imo.
    if t[i].Type == "item" {
      return true
    } else if t[j].Type == "item" {
      return false
    }
    return t[i].Type < t[j].Type
  } else {
    // sort by first argument
    return t[i].Args[0] < t[j].Args[0]
  }

  return false
}

func canDisplaySlot(slot int) bool {
  for _, i := range nameSlotMap {
    if slot == i {
      return true
    }
  }
  return false
}

func canDisplayName(name string) bool {
  for i, _ := range nameSlotMap {
    if name == i {
      return true
    }
  }
  return false
}

func main() {
  flag.Parse()

  if *versionFlag {
    fmt.Println(VERSION.String())
    return
  }

  url := flag.Arg(0)

  // filename, _ := osext.Executable()
  // APPDIR = filepath.Dir(filename)

  if url == "" {
    RunUI()
    return
  }

  err := Generate(map[string]interface{}{"url": url}, os.Stdout)
  if err != nil {
    io.WriteString(os.Stdout, err.Error())
  }
}

// options:
//  url - url to generate codes from
//  notmog - turn off grabbing transmogged items from armory
func Generate(options map[string]interface{}, w io.Writer) error {
  url := to.String(options["url"])
  var tmorphItems TMorphItems
  var err error
  switch {
  case strings.Contains(url, "wowhead.com"):
    tmorphItems, err = wowhead(options)
  case strings.Contains(url, "battle.net/wow"):
    tmorphItems, err = wowarmory(options)
  case strings.Contains(url, "http"):
    tmorphItems, err = generic(options)
  default:
    return errors.New("Do not recognize the URL.")
  }

  if err != nil {
    return err
  }

  sort.Sort(tmorphItems)
  for _, item := range tmorphItems {
    io.WriteString(w, item.String())
    io.WriteString(w, "\n")
  }

  return nil
}
