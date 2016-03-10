package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gosexy/to"
)

var (
	itemIdRe = regexp.MustCompile(`su_addToSaved\((.*?), (\d+)\)`)
)

func wowhead(options map[string]interface{}) (TMorphItems, error) {
	url := to.String(options["url"])

	// if they just put a wowhead item url in, just output that item
	if matches := wowheadUrlRe.FindStringSubmatch(url); len(matches) > 0 {
		items, err := wowapi([]string{matches[1]})
		if err != nil {
			return nil, errors.New(merry.Details(err))
		}

		if len(items) > 0 {
			return items, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	if matches := itemIdRe.FindStringSubmatch(string(data)); len(matches) == 3 {
		//return nil, errors.New(fmt.Sprintf("%#+v", matches))
		count, _ := strconv.Atoi(matches[2])
		if count == 1 {
			items, err := wowapi([]string{matches[1]})
			if err != nil {
				return nil, errors.New(merry.Details(err))
			}

			if len(items) > 0 {
				return items, nil
			}
		} else if count > 1 {
			itemids := strings.Split(matches[1][1:len(matches[1])-1], ":")
			items, err := wowapi(itemids)
			if err != nil {
				return nil, errors.New(merry.Details(err))
			}

			if len(items) > 0 {
				return items, nil
			}
		}
	}

	return nil, errors.New(`Could not find anything to morph on that wowhead page.`)
}

// package main
//
// // TODO: could use the wowhead item xml api.
// // e.g. http://www.wowhead.com/item=113939&xml
// // does not work for spells or itemsets...
//
// import (
// 	"bytes"
// 	"errors"
// 	"io/ioutil"
// 	"net/http"
// 	nurl "net/url"
// 	"regexp"
//
// 	"github.com/PuerkitoBio/goquery"
// 	"github.com/ansel1/merry"
// 	"github.com/gosexy/to"
// 	"github.com/robertkrimen/otto"
//
// 	"github.com/spf13/nitro"
// )
//
// var (
// 	itemRe          = regexp.MustCompile(`(?s)\<script type="text/javascript"\>//\<!\[CDATA\[.*?(g_items\.add.*?)//\]\]\>\</script\>`)
// 	wowheadNPCUrlRe = regexp.MustCompile(`wowhead.com/\??npc=(\d+)`)
// 	displayIdRe     = regexp.MustCompile(`displayId: (\d+)`)
// 	spellIdRe       = regexp.MustCompile(`<a href="(.*?)" class="q2">`)
//
// 	Timer *nitro.B
// )
//
// func init() {
// 	Timer = nitro.Initalize()
// 	nitro.AnalysisOn = false
// }
//
// const (
// 	// a helper that replaces some variables that wowhead would
// 	// have with my own objects. The main part is g_items where
// 	// I set the add function of that object to processItem, which
// 	// we will inject into the javascript using otto.Set().
// 	jscriptHelper = `
// g_items = {add: processItem};
// Summary = function(){};
// `
// )
//
// func fixWowheadURL(u string) (string, error) {
// 	tmpurl, err := nurl.Parse(u)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	if tmpurl.Host == "" {
// 		tmpurl.Host = "wowhead.com"
// 	}
//
// 	if tmpurl.Scheme == "" {
// 		tmpurl.Scheme = "http"
// 	}
//
// 	return tmpurl.String(), nil
// }
//
// // Parses the wowhead html and finds where it's adding items
// // to the comparison list via g_items.add. I use otto and a
// // little javascript helper and a javascript function to parse
// // and interpret the javascript and find all of the item data I need.
// // Example of the javascript we will parse:
// //   g_items.add(22423, {name_enus:'Dreadnaught Bracers', quality:4,icon:'INV_Bracer_15',jsonequip:{...}});
// func wowhead(options map[string]interface{}) (TMorphItems, error) {
// 	url := to.String(options["url"])
//
// 	Timer.Step("started")
//
// 	// if they just put a wowhead item url in, just output that item
// 	if matches := wowheadUrlRe.FindStringSubmatch(url); len(matches) > 0 {
// 		items, err := wowapi([]string{matches[1]})
// 		if err != nil {
// 			return nil, errors.New(merry.Details(err))
// 		}
//
// 		if len(items) > 0 {
// 			return items, nil
// 		}
// 	}
//
// 	Timer.Step("wowhead url match")
//
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	data, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp.Body.Close()
//
// 	Timer.Step("grab doc")
//
// 	if match := itemRe.FindSubmatch(data); len(match) > 0 {
// 		o := otto.New()
//
// 		// We fill these in to make it work for wowhead transmog sets.
// 		dollarObj, _ := o.Object(`$ = {}`)
// 		dollarObj.Set("extend", func(call otto.FunctionCall) otto.Value {
// 			return otto.UndefinedValue()
// 		})
// 		o.Set("$", dollarObj)
// 		g_spellsObj, _ := o.Object(`g_spells = {}`)
// 		o.Set("g_spells", g_spellsObj)
// 		o.Set("ts_PopulateScreenshotDiv", func(call otto.FunctionCall) otto.Value {
// 			return otto.UndefinedValue()
// 		})
//
// 		var tmorphItems TMorphItems
// 		seenMainHand := false
// 		// Our processItem function that gets called via the g_items.add() call
// 		// and the jscriptHelper script.
// 		o.Set("processItem", func(call otto.FunctionCall) otto.Value {
// 			// data we want is in the second argument
// 			v, _ := call.Argument(1).Export()
// 			// we're only interested in the jsonequip map
// 			datam := Map(v.(map[string]interface{})["jsonequip"])
//
// 			slot := int(to.Int64(datam["slot"]))
// 			id := int(to.Int64(datam["id"]))
// 			if v, ok := slotMap[slot]; ok {
// 				slot = v
// 			}
//
// 			if canDisplaySlot(slot) {
// 				// We're going to assume if someone has a list
// 				// that contains two main hands, they mean they want
// 				// it in their main hand and off hand.
// 				if slot == 16 {
// 					if seenMainHand {
// 						slot = 17
// 					}
// 					seenMainHand = true
// 				}
//
// 				tmorphItems = append(tmorphItems, &TMorphItem{
// 					Type: "item",
// 					Args: []int{slot, id},
// 				})
// 			}
//
// 			return otto.UndefinedValue()
// 		})
//
// 		// run the
// 		_, err = o.Run(jscriptHelper + string(match[1]))
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		if len(tmorphItems) > 0 {
// 			return tmorphItems, nil
// 		}
// 	}
//
// 	Timer.Step("find items in wowhead page")
//
// 	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
// 	if err != nil {
// 		return nil, errors.New("Could not parse wowhead NPC html")
// 	}
//
// 	Timer.Step("make goquery doc")
//
// 	// we got an npc. find the display id
// 	if matches := wowheadNPCUrlRe.FindStringSubmatch(url); len(matches) > 0 {
// 		node := doc.Find(`a:contains("View in 3D")`)
// 		if node.Length() == 0 {
// 			return nil, errors.New(`Unable to find "View in 3D" on page.`)
// 		}
//
// 		onclick, ok := node.Attr("onclick")
// 		if !ok {
// 			return nil, errors.New(`Unable to find "onclick" handler for "View in 3D" link`)
// 		}
//
// 		matches := displayIdRe.FindStringSubmatch(onclick)
// 		if len(matches) <= 1 {
// 			return nil, errors.New(`Unable to find display ID`)
// 		}
//
// 		return TMorphItems{
// 			&TMorphItem{
// 				Type: "morph",
// 				Args: []int{int(to.Int64(matches[1]))},
// 			},
// 		}, nil
// 	}
//
// 	Timer.Step("find npc display id")
//
// 	// try to find a url in the effect that we can parse for a displayid.
// 	node := doc.Find(`th:contains("Effect")`)
// 	if node.Length() > 0 {
// 		node = node.Next().Find("a")
// 		if node.Length() > 0 {
// 			if urltext, ok := node.Attr("href"); ok {
// 				fixedurl, err := fixWowheadURL(urltext)
// 				if err == nil {
// 					// just rerun the function with the found url
// 					return wowhead(map[string]interface{}{
// 						"url": fixedurl,
// 					})
// 				}
// 			}
// 		}
// 	}
//
// 	Timer.Step("find effect")
//
// 	// look for a spell id in a tooltip that may have a spell effect we
// 	// need to parse
// 	if matches := spellIdRe.FindStringSubmatch(string(data)); len(matches) > 1 {
// 		fixedurl, err := fixWowheadURL(matches[1])
// 		if err == nil {
// 			// rerun the function with the url we found
// 			return wowhead(map[string]interface{}{
// 				"url": fixedurl,
// 			})
// 		}
// 	}
//
// 	Timer.Step("find spell id")
//
// 	return nil, errors.New(`Could not find anything to morph on that wowhead page.`)
// }
