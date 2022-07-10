package xconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type setter func(value interface{}) error

type appSetter struct {
	appId       string
	fieldSetter map[string]setter
}

func (as *appSetter) Reset() error {
	for key, f := range as.fieldSetter {
		if err := f(getLocalConfig(as.appId, key)); err != nil {
			return err
		}
	}

	return nil
}

func (as *appSetter) AddSetter(key string, s setter) {
	as.fieldSetter[key] = s
}

func (as *appSetter) SetValue(key string, v interface{}) error {
	if f, ok := as.fieldSetter[key]; ok {
		if hasLocalConfig(as.appId, key) {
			return nil
		}

		return f(v)
	}

	return nil
}

type parser struct {
	appSetter map[string]*appSetter
}

func newParser() *parser {
	return &parser{appSetter: map[string]*appSetter{}}
}

func (p *parser) parse(obj interface{}) (map[string]*appSetter, error) {
	rv, err := p.validateReflectValue(obj)
	if err != nil {
		return nil, err
	}

	if err := p.packup(rv, tag{}, "", ""); err != nil {
		return nil, err
	}

	return p.appSetter, nil
}

func (p *parser) validateReflectValue(obj interface{}) (rv reflect.Value, err error) {
	if obj == nil {
		err = errors.New("no data provided")
		return
	}

	vv := reflect.ValueOf(obj)
	switch vv.Kind() {
	case reflect.Pointer:
		if vv.IsNil() {
			err = errors.New("nil pointer of a struct is not supported")
			return
		}

		vv = vv.Elem()
		if vv.Kind() != reflect.Struct {
			err = fmt.Errorf("invalid type of input : %v", vv.Type())
			return
		}
		rv = vv
	default:
		err = fmt.Errorf("invalid type of input : %v", vv.Type())
	}

	return
}

//将结构体中的字段做映射
func (p *parser) packup(rv reflect.Value, t tag, ns, appId string) error {
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		if rv.Elem().Kind() == reflect.Pointer {
			return nil
		}

		return p.packup(rv.Elem(), t, ns, appId)
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			ft := rv.Type().Field(i)

			n := ns
			if !ft.Anonymous {
				t = parseTag(rv.Type().Field(i))
				if t.IsIgnore() {
					continue
				}
				n = parseNs(n, t)
				if t.appId != "" {
					appId = t.appId
				}
			}

			rvv := rv.Field(i)
			if err := p.packup(rvv, t, n, appId); err != nil {
				return err
			}
		}
	case reflect.String:
		p.addSetter(appId, ns, p.setString(rv, t.defaultValue))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.addSetter(appId, ns, p.setInt(rv, t.defaultValue))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		p.addSetter(appId, ns, p.setUint(rv, t.defaultValue))
	case reflect.Bool:
		p.addSetter(appId, ns, p.setBool(rv, t.defaultValue))
	}

	return nil
}

func (p *parser) addSetter(appId, ns string, f setter) {
	var set *appSetter
	set, ok := p.appSetter[appId]
	if !ok {
		set = &appSetter{
			appId:       appId,
			fieldSetter: map[string]setter{},
		}
	}

	set.AddSetter(ns, f)
	p.appSetter[appId] = set
}

func (p *parser) setString(rv reflect.Value, defaultValue interface{}) setter {
	return func(value interface{}) error {
		if value == nil {
			value = defaultValue
		}

		rv.Set(reflect.ValueOf(value))

		return nil
	}
}

func (p *parser) setInt(rv reflect.Value, defaultValue interface{}) setter {
	return func(value interface{}) (err error) {
		if value == nil {
			value = defaultValue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			iv, err := strconv.ParseInt(value.(string), 10, 64)
			if err != nil {
				return err
			}
			rv.Set(reflect.ValueOf(iv).Convert(rv.Type()))
		default:
			rv.Set(reflect.ValueOf(value).Convert(rv.Type()))
		}

		return nil
	}
}

func (p *parser) setUint(rv reflect.Value, defaultValue interface{}) setter {
	return func(value interface{}) error {
		if value == nil {
			value = defaultValue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			iv, err := strconv.ParseUint(value.(string), 10, 64)
			if err != nil {
				return err
			}

			rv.Set(reflect.ValueOf(iv).Convert(rv.Type()))
		default:
			rv.Set(reflect.ValueOf(value).Convert(rv.Type()))
		}

		return nil
	}
}

func (p *parser) setBool(rv reflect.Value, defaultValue interface{}) setter {
	return func(value interface{}) error {
		if value == nil {
			value = defaultValue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			bv, err := strconv.ParseBool(value.(string))
			if err != nil {
				return err
			}

			rv.Set(reflect.ValueOf(bv))
		default:
			rv.Set(reflect.ValueOf(value))
		}

		return nil
	}
}

func parseNs(ns string, t tag) string {
	build := strings.Builder{}
	build.WriteString(ns)
	if ns != "" {
		build.WriteByte('.')
	}
	build.WriteString(t.name)

	return build.String()
}
