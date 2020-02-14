package berus

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type Custom interface{}
type custom struct{ Custom }

var cSlice = make(map[string]Custom)

func RegisterCustom(name string, c Custom) {
	cSlice[name] = c
}

func CustomHookFunc(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {
	if f.Kind() != reflect.Map || t.Kind() != reflect.Interface || !reflect.TypeOf((*custom)(nil)).Implements(t){
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
	return c, mapstructure.Decode(val, c)
}
