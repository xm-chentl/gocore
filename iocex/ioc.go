package iocex

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	tagInject  = "inject"
	TagDefault = "default"
)

var container sync.Map

func Get(inst interface{}) interface{} {
	rt := getType(inst)
	if v, ok := container.Load(rt); ok {
		return v
	}

	return nil
}

func getTag(inst interface{}, tagName string) interface{} {
	rt := getType(inst)
	v, ok := container.Load(rt)
	if !ok {
		return nil
	}

	vMap, ok := v.(map[string]interface{})
	if !ok {
		return v
	}
	if v, ok = vMap[tagName]; ok {
		return v
	}

	return nil
}

func Has(inst interface{}) (ok bool) {
	_, ok = container.Load(getType(inst))
	return
}

func Inject(structInst interface{}, funcs ...func(reflect.StructField) interface{}) (err error) {
	rt := reflect.TypeOf(structInst)
	rv := reflect.ValueOf(structInst)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		originalRt := field.Type
		fieldRt := field.Type
		if fieldRt.Kind() == reflect.Ptr {
			fieldRt = fieldRt.Elem()
		}
		if len(funcs) > 0 && funcs[0] != nil {
			v := funcs[0](field)
			if v != nil {
				rv.Field(i).Set(reflect.ValueOf(v))
				continue
			}
		}
		if fieldRt.Kind() == reflect.Interface {
			name, ok := field.Tag.Lookup(tagInject)
			if !ok {
				continue
			}

			v, ok := container.Load(field.Type)
			if !ok {
				continue
			}

			vv := reflect.ValueOf(v)
			if vv.Kind() == reflect.Map {
				if name == "" {
					name = TagDefault
				}

				var fieldRv reflect.Value
				mapIter := vv.MapRange()
				for mapIter.Next() {
					if strings.EqualFold(mapIter.Key().String(), name) {
						fieldRv = mapIter.Value()
						break
					}
				}
				// todo: 缺少默认配置的问题
				if fieldRv.IsValid() {
					v = fieldRv.Interface()
				} else {
					err = fmt.Errorf("specifies that the %s is not registered(%s)", rt.Name(), field.Type.Name())
					return
				}
			}
			if v != nil {
				rv.Field(i).Set(reflect.ValueOf(v))
			}
		} else if fieldRt.Kind() == reflect.Struct {
			inst := reflect.New(fieldRt).Interface()
			err = Inject(inst)
			if err != nil {
				return
			}
			if inst != nil {
				if originalRt.Kind() == reflect.Ptr {
					rv.Field(i).Set(reflect.ValueOf(inst))
				} else {
					rv.Field(i).Set(reflect.ValueOf(inst).Elem())
				}
			}
		}

	}
	return
}

func Set(key interface{}, inst interface{}) {
	rt := getType(key)
	container.Store(rt, inst)
}

func SetMap(key interface{}, mapping map[string]interface{}) {
	rt := getType(key)
	container.Store(rt, mapping)
}

func getType(v interface{}) reflect.Type {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Interface {
		panic("Specifies that the object parameter is not an interface")
	}

	return rt
}
