package main

import (
	"github.com/gosexy/to"
	"reflect"
)

func Map(val interface{}) map[string]interface{} {
	list := map[string]interface{}{}
	if val == nil {
		return list
	}
	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		vval := reflect.ValueOf(val)
		vlist := reflect.ValueOf(list)
		for _, vkey := range vval.MapKeys() {
			key := to.String(vkey.Interface())
			vlist.SetMapIndex(reflect.ValueOf(key), vval.MapIndex(vkey))
		}
		return list
	}
	return list
}
