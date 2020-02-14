package berus

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type Custom interface{}

var cSlice = make(map[string]Custom)
var empty Custom

func RegisterCustom(name string, c Custom) {
	cSlice[name] = c
}

func CustomHookFunc(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {

	if f.Kind() != reflect.Map || t.Kind() != reflect.Interface {
		return data, nil
	}
	if of := reflect.TypeOf(&empty); of.Elem() == nil || !of.Elem().Implements(t) {
		return data, nil
	}
	val, ok := data.(map[string]interface{})
	if !ok {
		return nil, newError("Unsupported data")
	}
	typ, ok := val["_type"]
	if !ok {
		return nil, newError("Custom doesn't have '_type'")
	}
	tt, ok := typ.(string)
	if !ok {
		return nil, newError("Unsupported field '_type'")
	}
	c, ok := cSlice[tt]
	if !ok {
		return nil, newError("Unregistered custom type")
	}
	delete(val, "_type")
	value := reflect.New(reflect.TypeOf(c).Elem()).Interface()
	return value, mapstructure.Decode(val, value)
}
