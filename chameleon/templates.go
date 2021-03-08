package chameleon

import (
	"fmt"
	"reflect"
	"text/template"
)

var funcs = template.FuncMap{
	"first": func(item reflect.Value) (reflect.Value, error) {
		item, isNil := indirect(item)
		if isNil {
			return reflect.Value{}, fmt.Errorf("index of nil pointer")
		}
		item = item.Index(0)
		return item, nil
	},
	"last": func(item reflect.Value) (reflect.Value, error) {
		item, isNil := indirect(item)
		if isNil {
			return reflect.Value{}, fmt.Errorf("index of nil pointer")
		}
		item = item.Index(item.Len() - 1)
		return item, nil
	},
}

func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}
