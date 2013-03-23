package main

import (
  "code.google.com/p/go-semver/version"
  // "io/ioutil"
  // "os"
  // "path/filepath"
)

var (
  VERSION = &version.Version{
    Major: 1,
    Minor: 0,
    Patch: 2,
  }
)

// type Version struct {
//   *version.Version
//   invalid bool
// }

// func (v *Version) String() string {
//   if v.invalid {
//     return "unknown"
//   }

//   return v.Version.String()
// }

// func (v *Version) SetInvalid() {
//   v.invalid = true
// }

// var (
//   VERSION *Version = &Version{}
// )

// func initVersion() {
//   file, err := os.Open(filepath.Join(APPDIR, "version.txt"))
//   if err != nil {
//     VERSION.SetInvalid()
//     return
//   }
//   defer file.Close()

//   data, err := ioutil.ReadAll(file)
//   if err != nil {
//     VERSION.SetInvalid()
//     return
//   }

//   v, err := version.Parse(string(data))
//   if err != nil {
//     VERSION.SetInvalid()
//     return
//   }

//   VERSION.Version = v
// }
