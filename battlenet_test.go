package main

import (
	"testing"

	"github.com/ansel1/merry"
)

func TestWowarmory(t *testing.T) {
	items, err := wowarmory(map[string]interface{}{
		"url":    "http://us.battle.net/wow/en/character/tichondrius/Nahj/simple",
		"apikey": "g6pnns5wzvy9zqb7vtduvqjente6yqrx",
	})
	if err != nil {
		t.Fatal(merry.Details(err))
	}
	if len(items) == 0 {
		t.Fatal("could not retrieve wowarmory items")
	}
}
