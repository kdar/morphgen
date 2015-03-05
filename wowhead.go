package main

import (
	"errors"
	"github.com/gosexy/to"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	itemRe = regexp.MustCompile(`(?s)\<script type="text/javascript"\>//\<!\[CDATA\[.*?(g_items\.add.*?)//\]\]\>\</script\>`)
)

const (
	// a helper that replaces some variables that wowhead would
	// have with my own objects. The main part is g_items where
	// I set the add function of that object to processItem, which
	// we will inject into the javascript using otto.Set().
	jscriptHelper = `
g_items = {add: processItem};
Summary = function(){};
`
)

// Parses the wowhead html and finds where it's adding items
// to the comparison list via g_items.add. I use otto and a
// little javascript helper and a javascript function to parse
// and interpret the javascript and find all of the item data I need.
// Example of the javascript we will parse:
//   g_items.add(22423, {name_enus:'Dreadnaught Bracers', quality:4,icon:'INV_Bracer_15',jsonequip:{...}});
func wowhead(options map[string]interface{}) (TMorphItems, error) {
	url := to.String(options["url"])

	// if they just put a wowhead item url in, just output that item
	if matches := wowheadUrlRe.FindStringSubmatch(url); len(matches) > 0 {
		return wowapi([]string{matches[1]})
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

	match := itemRe.FindSubmatch(data)
	if len(match) == 0 {
		return nil, errors.New("Could not find wowhead items")
	}

	o := otto.New()

	// We fill these in to make it work for wowhead transmog sets.
	dollarObj, _ := o.Object(`$ = {}`)
	dollarObj.Set("extend", func(call otto.FunctionCall) otto.Value {
		return otto.UndefinedValue()
	})
	o.Set("$", dollarObj)
	g_spellsObj, _ := o.Object(`g_spells = {}`)
	o.Set("g_spells", g_spellsObj)

	var tmorphItems TMorphItems
	seenMainHand := false
	// Our processItem function that gets called via the g_items.add() call
	// and the jscriptHelper script.
	o.Set("processItem", func(call otto.FunctionCall) otto.Value {
		// data we want is in the second argument
		v, _ := call.Argument(1).Export()
		// we're only interested in the jsonequip map
		datam := Map(v.(map[string]interface{})["jsonequip"])

		slot := int(to.Int64(datam["slot"]))
		id := int(to.Int64(datam["id"]))
		if v, ok := slotMap[slot]; ok {
			slot = v
		}

		if canDisplaySlot(slot) {
			// We're going to assume if someone has a list
			// that contains two main hands, they mean they want
			// it in their main hand and off hand.
			if slot == 16 {
				if seenMainHand {
					slot = 17
				}
				seenMainHand = true
			}

			tmorphItems = append(tmorphItems, &TMorphItem{
				Type: "item",
				Args: []int{slot, id},
			})
		}

		return otto.UndefinedValue()
	})

	// run the
	_, err = o.Run(jscriptHelper + string(match[1]))
	if err != nil {
		return nil, err
	}

	return tmorphItems, nil
}
