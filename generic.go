package main

import (
  "github.com/gosexy/to"
  "io/ioutil"
  "net/http"
  "regexp"
  //"strings"
)

var (
  wowheadUrlRe = regexp.MustCompile(`wowhead.com/\??item=(\d+)`)
  wowheadUrl   = "http://www.wowhead.com/compare?items="
)

// Attempt to find any links in the page that we can
// parse and generate codes for.
func generic(options map[string]interface{}) (TMorphItems, error) {
  resp, err := http.Get(to.String(options["url"]))
  if err != nil {
    return nil, err
  }

  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  resp.Body.Close()

  matches := wowheadUrlRe.FindAllStringSubmatch(string(data), -1)
  var items []string
  for _, match := range matches {
    items = append(items, match[1])
  }

  //options["url"] = wowheadUrl + strings.Join(items, ";")
  //return wowhead(options)
  return wowapi(items)
}
