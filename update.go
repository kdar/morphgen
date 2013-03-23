package main

import (
  "code.google.com/p/go-semver/version"
  "encoding/json"
  "fmt"
  "github.com/gosexy/to"
  "io/ioutil"
  "net/http"
  "os/exec"
  "runtime"
)

// func CheckUpdate() (string, error) {
//   resp, err := http.Get("https://api.github.com/repos/kdar/morphgen-binary/contents/version.txt")
//   if err != nil {
//     return "", err
//   }

//   data, err := ioutil.ReadAll(resp.Body)
//   if err != nil {
//     return "", err
//   }
//   resp.Body.Close()

//   var datam map[string]interface{}
//   err = json.Unmarshal(data, &datam)
//   if err != nil {
//     return "", err
//   }

//   b64version := strings.Trim(to.String(datam["content"]), "\n")
//   versionData, err := base64.StdEncoding.DecodeString(b64version)
//   if err != nil {
//     return "", err
//   }

//   v, err := version.Parse(string(versionData))
//   if err != nil {
//     return "", err
//   }

//   if VERSION.Less(v) {
//     return "Update available: " + v.String(), nil
//   }

//   return "", nil
// }

func CheckUpdate() (string, error) {
  resp, err := http.Get("https://api.github.com/repos/kdar/morphgen-binary/tags")
  if err != nil {
    return "", err
  }

  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }
  resp.Body.Close()

  var datam []map[string]interface{}
  err = json.Unmarshal(data, &datam)
  if err != nil {
    return "", err
  }

  if len(datam) > 0 {
    v, err := version.Parse(to.String(datam[0]["name"]))
    if err != nil {
      return "", err
    }

    if VERSION.Less(v) {
      return "Update available: " + v.String(), nil
    }
  }

  return "", nil
}

func OpenDownloadInBrowser() error {
  url := "https://github.com/kdar/morphgen-binary/archive/master.zip"

  var err error
  switch runtime.GOOS {
  case "linux":
    err = exec.Command("xdg-open", url).Start()
  case "windows":
    err = exec.Command("open", url).Start()
    if err != nil {
      err = exec.Command(`rundll32.exe`, "url.dll,FileProtocolHandler", url).Start()
      if err != nil {
        err = exec.Command(`C:\Windows\System32\rundll32.exe`, "url.dll,FileProtocolHandler", url).Start()
      }
    }
  case "darwin":
    err = exec.Command("open", url).Start()
  default:
    err = fmt.Errorf("unsupported platform")
  }

  return err
}
