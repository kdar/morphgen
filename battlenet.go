package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gosexy/to"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	// "path"
	"strings"
)

var (
	appearanceMap = map[string]string{
		"faceVariation": "face",
		"skinColor":     "skin",
		"hairVariation": "hair",
		"hairColor":     "haircolor",
		//"featureVariation": ,
		//"showHelm": ,
		//"showCloak": ,
	}
)

// Get detailed codes from a character's armory page.
func wowarmory(options map[string]interface{}) (TMorphItems, error) {
	u, err := neturl.Parse(to.String(options["url"]))
	if err != nil {
		return nil, err
	}

	parts := strings.Split(u.Path, "/")
	loc := -1
	for x := 0; x < len(parts); x++ {
		if parts[x] == "character" {
			loc = x + 1
			break
		}
	}

	if loc == -1 {
		return nil, errors.New("Could not parse battle.net URL")
	}

	u.Path = fmt.Sprintf("/api/wow/character/%s/%s", parts[loc], parts[loc+1])
	u.RawQuery = "fields=items,appearance"

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	var tmorphItems TMorphItems
	var datam map[string]interface{}

	err = json.Unmarshal(data, &datam)
	if err != nil {
		return nil, err
	}

	// get all armor, weapons, and enchants
	items := Map(datam["items"])
	for k, v := range items {
		if v, ok := v.(map[string]interface{}); ok {
			if canDisplayName(k) {
				id := to.Int64(v["id"])
				tooltipParams := Map(v["tooltipParams"])
				if !to.Bool(options["notmog"]) {
					transmogItem := tooltipParams["transmogItem"]
					if transmogItem != nil {
						id = to.Int64(transmogItem)
					}
				}

				tmorphItems = append(tmorphItems, &TMorphItem{
					Type: "item",
					Args: []int{nameSlotMap[k], int(id)},
				})

				// get enchants off the weapons
				if k == "mainHand" || k == "offHand" {
					if tooltipParams["enchant"] != nil {
						which := 1
						if k == "offHand" {
							which = 2
						}

						tmorphItems = append(tmorphItems, &TMorphItem{
							Type: "enchant",
							Args: []int{which, int(to.Int64(tooltipParams["enchant"]))},
						})
					}
				}
			}
		}
	}

	// set offhand to 0 if there is none.
	// TODO: maybe there should be defaults?
	if items["offHand"] == nil {
		tmorphItems = append(tmorphItems, &TMorphItem{
			Type: "item",
			Args: []int{nameSlotMap["offHand"], 0},
		})
	}

	// appearance stuff
	appearance := Map(datam["appearance"])
	for k, v := range appearance {
		if typ, ok := appearanceMap[k]; ok {
			tmorphItems = append(tmorphItems, &TMorphItem{
				Type: typ,
				Args: []int{int(to.Int64(v))},
			})
		}
	}
	tmorphItems = append(tmorphItems, &TMorphItem{
		Type: "race",
		Args: []int{int(to.Int64(datam["race"]))},
	}, &TMorphItem{
		Type: "gender",
		Args: []int{int(to.Int64(datam["gender"]))},
	})

	return tmorphItems, nil
}

// Get the codes for the list of ids via the wow api.
func wowapi(ids []string) (TMorphItems, error) {
	var tmorphItems TMorphItems
	idslen := len(ids)
	errChan := make(chan error)
	doneChan := make(chan bool, idslen)

	for _, id := range ids {
		go func(id string) {
			defer func() {
				doneChan <- true
			}()

			resp, err := http.Get("http://us.battle.net/api/wow/item/" + id)
			if err != nil {
				errChan <- err
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errChan <- err
				return
			}
			resp.Body.Close()

			var datam map[string]interface{}
			err = json.Unmarshal(data, &datam)
			if err != nil {
				errChan <- err
				return
			}

			slot := int(to.Int64(datam["inventoryType"]))
			if v, ok := slotMap[slot]; ok {
				slot = v
			}
			if canDisplaySlot(slot) {
				tmorphItems = append(tmorphItems, &TMorphItem{
					Type: "item",
					Args: []int{slot, int(to.Int64(id))},
				})
			}
		}(id)
	}

	count := 0
	for count < idslen {
		select {
		case err := <-errChan:
			return nil, err
		case <-doneChan:
			count++
			if count >= idslen {
				return tmorphItems, nil
			}
		}
	}

	return nil, nil
}
