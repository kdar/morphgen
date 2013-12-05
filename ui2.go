package main

import (
  "bytes"
  "code.google.com/p/go.net/websocket"
  "fmt"
  "github.com/gosexy/to"
  "math/rand"
  "net"
  "net/http"
  "os"
  "os/exec"
  "path/filepath"
  "time"
)

var (
  isConnected = false
  port        = 0
  CHROME_PATH = ""
)

type Message struct {
  Name string
  Args interface{}
}

func errToJs(err error) interface{} {
  if err != nil {
    return err.Error()
  }
  return nil
}

func Wgenerate(ws *websocket.Conn, args interface{}) {
  go func() {
    buffer := &bytes.Buffer{}
    err := Generate(to.Map(args), buffer)

    args := []interface{}{buffer.String(), errToJs(err)}
    websocket.JSON.Send(ws, Message{"generate_callback", args})
  }()
}

func Wcheckupdate(ws *websocket.Conn, args interface{}) {
  go func() {
    update, err := CheckUpdate()
    args := []interface{}{update, errToJs(err)}
    websocket.JSON.Send(ws, Message{"checkupdate_callback", args})
  }()
}

func wsHandler(ws *websocket.Conn) {
  isConnected = true
  for {
    var msg Message
    // Receive receives a text message serialized T as JSON.
    err := websocket.JSON.Receive(ws, &msg)
    if err != nil {
      break
    }

    switch msg.Name {
    case "checkupdate":
      Wcheckupdate(ws, msg.Args)
    case "generate":
      Wgenerate(ws, msg.Args)
    }
  }
  os.Exit(0)
}

func runChrome() {
  // cache bust the url
  // url := fmt.Sprintf("http://localhost:%d/?cb=%d", port, rand.Int())
  url := fmt.Sprintf("http://localhost:%d/", port)
  // `--user-data-dir=c:\cacaface`
  cmd := exec.Command(CHROME_PATH, fmt.Sprintf(`--app=%s`, url))
  err := cmd.Start()
  if err != nil {
    panic(err)
  }
}

func RunUI2() {
  rand.Seed(time.Now().UnixNano())

  http.Handle("/ws", websocket.Handler(wsHandler))
  http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./public"))))

  ch := time.After(15 * time.Second)
  go func() {
    <-ch
    if !isConnected {
      fmt.Println("GUI never connected to backend.")
      os.Exit(-1)
    }
  }()

  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    panic(err)
  }

  port = l.Addr().(*net.TCPAddr).Port
  runChrome()

  http.Serve(l, nil)
}

func init() {
  // err := wax.SHGetFolderPath(0, wax.CSIDL_LOCAL_APPDATA, 0, 0, &APP_DATA)
  // if err == nil {
  //   CHROME_PATH = APP_DATA + "\\Google\\Chrome\\Application\\chrome.exe"
  //   DATA_PATH = APP_DATA + "\\"
  // } else {
  //   MessageBoxFatal(err.Error(), "Error")
  // }

  // VERY slow for some reason
  // u, err := user.Current()
  // if err != nil {
  //   panic(err)
  // }

  appData := os.Getenv("LOCALAPPDATA")

  CHROME_PATH = filepath.Join(appData, "Google", "Chrome", "Application", "chrome.exe")
}
