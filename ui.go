package main

import (
  "bytes"
  "github.com/kdar/morphgen/golua/lua"
  "github.com/kdar/morphgen/luar"
  "reflect"
)

// Add a callback to the lua state. 
func addCallback(L *lua.State, name string, args []interface{}) {
  L.GetGlobal("callbacks")
  luar.GoToLua(L, reflect.TypeOf(args), reflect.ValueOf(args))
  L.SetField(-2, name)
}

// A lua function wrapper to Generate().
func Lgenerate(L *lua.State) int {
  v := luar.CopyTableToMap(L, nil, 1).(map[string]interface{})

  go func() {
    buffer := &bytes.Buffer{}
    err := Generate(v, buffer)

    args := []interface{}{buffer.String(), err}
    addCallback(L, "generate_callback", args)
  }()
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
  L.Register("checkupdate", Lcheckupdate)

  luar.Register(L, "", luar.Map{
    "download": Ldownload,
  })

  err := L.DoFile("ui/ui.lua")
  if err != nil {
    panic(err)
  }
}
