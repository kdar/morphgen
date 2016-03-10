package main

import (
	// "bitbucket.org/kardianos/osext"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gosexy/to"
	// "path/filepath"
	"sort"
	"strings"
)

var (
	APPDIR string

	generator *Generator

	noTmog      = flag.Bool("notmog", false, "Toggle off grabbing transmogged items from armory")
	versionFlag = flag.Bool("version", false, "Display version")

	// a map from names to slots. this also contains
	// what we're able to tmorph into
	nameSlotMap = map[string]int{
		"head":     1,
		"neck":     2,
		"shoulder": 3,
		"shirt":    4,
		"chest":    5,
		"waist":    6,
		"legs":     7,
		"feet":     8,
		"wrist":    9,
		"hands":    10,
		"finger":   11,
		"trinket":  12,
		"one-hand": 13,
		"shield":   14,
		"back":     15,
		"mainHand": 16,
		"offHand":  17,
		"bag":      18,
		"tabard":   19,
		"robe":     20,
	}

	// mappings of items type ids to wow slot ids
	slotMap = map[int]int{
		20: 5, // robe -> chest,

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

const (
	raidInherit = -1
	raidNormal  = 0
	raidHeroic  = 566
	raidMythic  = 567
)

type TMorphItem struct {
	Type  string
	Args  []int
	Bonus int
}

func (t TMorphItem) String() string {
	var args []string
	for _, i := range t.Args {
		args = append(args, to.String(i))
	}

	if t.Type == "item" {
		switch t.Bonus {
		case raidNormal:
			args = append(args, "0")
		case raidHeroic:
			args = append(args, "1")
		case raidMythic:
			args = append(args, "3")
		}
	}

	return fmt.Sprintf(".%s %s", t.Type, strings.Join(args, " "))
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
	generator = &Generator{}

	flag.Parse()

	if *versionFlag {
		fmt.Println(VERSION.String())
		return
	}

	url := flag.Arg(0)

	// filename, _ := osext.Executable()
	// APPDIR = filepath.Dir(filename)

	if url == "" {
		// RunUI2()
		RunUI()
		return
	}

	err := generator.Generate(map[string]interface{}{"url": url}, os.Stdout)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
}

type Generator struct {
	lastTmorphItems TMorphItems
}

// options:
//  url - url to generate codes from
//  notmog - turn off grabbing transmogged items from armory
func (g *Generator) Generate(options map[string]interface{}, w io.Writer) error {
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

	g.lastTmorphItems = tmorphItems
	bonus := int(to.Int64(options["bonus"]))
	g.Bonus(bonus)
	g.Output(w)

	return nil
}

func (g *Generator) Bonus(bonus int) {
	if bonus == raidInherit {
		return
	}

	for i, _ := range g.lastTmorphItems {
		g.lastTmorphItems[i].Bonus = bonus
	}
}

func (g *Generator) Output(w io.Writer) {
	sort.Sort(g.lastTmorphItems)
	for _, item := range g.lastTmorphItems {
		io.WriteString(w, item.String())
		io.WriteString(w, "\n")
	}
}
