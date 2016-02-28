package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"

	"github.com/ansel1/merry"
	"github.com/gosexy/to"
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

func apicall(u *neturl.URL) (*http.Response, error) {
	apikeys := []string{
		"sehm98r8ss3fddjfs8rwxur435xjav94",
		"trtrn4tdhzeruuypctm55a8m5td9z4ce",
		"cc7qrdaxn9h47ds39nmgzb9hnutxfsrr",
		"cbygdtw624xmfn2hr2sw3ayaeekg9cvy",
		"7zsnd2q7npg5bch8hacae7t28xskn8vr",
	}

	// if val, ok := options["apikey"]; ok {
	// 	apikeys = append([]string{to.String(val)}, apikeys...)
	// }

	rawQuery := u.RawQuery
	if len(rawQuery) > 0 {
		rawQuery += "&apikey="
	} else {
		rawQuery += "apikey="
	}

	var resp *http.Response
	var apierr error
	for _, apikey := range apikeys {
		u.RawQuery = rawQuery + apikey
		resp, apierr = http.Get(u.String())
		if apierr == nil {
			return resp, nil
		}
	}

	return nil, merry.Wrap(apierr)
}

// Get detailed codes from a character's armory page.
func wowarmory(options map[string]interface{}) (TMorphItems, error) {
	u, err := neturl.Parse(to.String(options["url"]))
	if err != nil {
		return nil, merry.Wrap(err)
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

	// FIXME: this isn't exactly correct, because you can be in the US and want
	// to get the tmorph codes of a person in EU/China. So we need to probably
	// have settings to where the user of TMorphGen is.
	hostParts := strings.Split(u.Host, ".")
	if hostParts[0] == "cn" {
	} else if len(hostParts) == 2 {
		u.Host = "us.api." + strings.Join(hostParts, ".")
	} else {
		u.Host = hostParts[0] + ".api." + strings.Join(hostParts[1:], ".")
	}

	u.Scheme = "https"
	u.Path = fmt.Sprintf("/wow/character/%s/%s", parts[loc], parts[loc+1])
	u.RawQuery = "fields=items,appearance&locale=en_US"

	resp, err := apicall(u)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, merry.Wrap(err)
	}
	resp.Body.Close()

	var tmorphItems TMorphItems
	var datam map[string]interface{}

	err = json.Unmarshal(data, &datam)
	if err != nil {
		return nil, merry.Wrap(err)
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

			contextUrl := ""

		REDO:
			// FIXME: using US api, ignoring user's location
			u, err := neturl.Parse("https://us.api.battle.net/wow/item/" + id + contextUrl)
			if err != nil {
				errChan <- err
				return
			}

			resp, err := apicall(u)
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

			if _, ok := datam["inventoryType"]; !ok {
				if contexts, ok := datam["availableContexts"]; ok {
					ctxs := contexts.([]interface{})
					contextUrl = "/" + ctxs[0].(string)
					goto REDO
				}
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
			return nil, merry.Wrap(err)
		case <-doneChan:
			count++
			if count >= idslen {
				return tmorphItems, nil
			}
		}
	}

	return nil, nil
}
