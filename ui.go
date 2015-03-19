package main

import (
	"bytes"
	"github.com/kdar/morphgen/golua/lua"
	"github.com/kdar/morphgen/luar"
	"reflect"
)

func bonusTextValue(t string) int {
	switch t {
	case "Inherit":
	case "Normal":
		return raidNormal
	case "Heroic":
		return raidHeroic
	case "Mythic":
		return raidMythic
	}

	return raidInherit
}

// Add a callback to the lua state.
func addCallback(L *lua.State, name string, args []interface{}) {
	L.GetGlobal("callbacks")
	luar.GoToLua(L, reflect.TypeOf(args), reflect.ValueOf(args))
	L.SetField(-2, name)
}

// A lua function wrapper to Generate().
func Lgenerate(L *lua.State) int {
	v := luar.CopyTableToMap(L, nil, 1).(map[string]interface{})
	v["bonus"] = bonusTextValue(v["bonustext"].(string))

	go func() {
		buffer := &bytes.Buffer{}
		err := generator.Generate(v, buffer)

		args := []interface{}{buffer.String(), err}
		addCallback(L, "generate_callback", args)
	}()
	return 0
}

func LonBonusChange(L *lua.State) int {
	v := luar.CopyTableToMap(L, nil, 1).(map[string]interface{})

	buffer := &bytes.Buffer{}

	generator.Bonus(bonusTextValue(v["text"].(string)))
	generator.Output(buffer)

	args := []interface{}{buffer.String(), nil}
	addCallback(L, "onBonusChange_callback", args)
	return 0
}

// A lua function wrapper to OpenDownloadInBrowser().
func Ldownload() error {
	return OpenDownloadInBrowser()
}

func Lcheckupdate(L *lua.State) int {
	go func() {
		update, err := CheckUpdate()
		args := []interface{}{update, err}
		addCallback(L, "checkupdate_callback", args)
	}()

	return 0
}

func RunUI() {
	L := luar.Init()
	defer L.Close()

	L.PushString(VERSION.String())
	L.SetGlobal("VERSION")

	L.Register("generate", Lgenerate)
	L.Register("onBonusChange", LonBonusChange)
	L.Register("checkupdate", Lcheckupdate)

	luar.Register(L, "", luar.Map{
		"download": Ldownload,
	})

	err := L.DoFile("ui/ui.lua")
	if err != nil {
		panic(err)
	}
}
