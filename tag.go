package xconfig

import (
	"reflect"
	"strings"
)

type tag struct {
	name         string
	defaultValue string
	ignore       bool
	appId        string
}

func (t tag) IsIgnore() bool {
	return t.ignore
}

func parseTag(sf reflect.StructField) tag {
	t := sf.Tag.Get("xconfig")
	if t == "-" {
		return tag{ignore: true}
	}

	tt := tag{name: genDefaultName(sf.Name), defaultValue: genDefaultValue(sf.Type.Kind())}

	fields := strings.Split(t, ";")
	for _, field := range fields {
		k, v := tt.parseTagField(field)
		switch k {
		case "name":
			tt.name = v
		case "default":
			tt.defaultValue = v
		case "appId":
			tt.appId = v
		}
	}

	return tt
}

//HelloWorld => helloWorld
func genDefaultName(a string) string {
	b := []byte(a)
	if len(b) > 0 && b[0] >= 'A' && b[0] <= 'Z' {
		b[0] = b[0] + 32
	}

	return string(b)
}

func genDefaultValue(kind reflect.Kind) (dv string) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dv = "0"
	case reflect.Bool:
		dv = "false"
	default:
		dv = ""
	}

	return
}

func (tag) parseTagField(filed string) (key, value string) {
	for i, c := range filed {
		if c == ':' {
			key = filed[:i]
			if i < len(filed) {
				value = filed[i+1:]
			}
		}
	}

	return
}
