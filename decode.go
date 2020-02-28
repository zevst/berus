package berus

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type Custom interface{}

var cSlice = make(map[string]Custom)
var empty Custom

// RegisterCustom. You must register a structure.
// After parsing, you can cast the field to this structure.
func RegisterCustom(name string, c Custom) {
	cSlice[name] = c
}

// CustomHookFunc decode hook for github.com/mitchellh/mapstructure
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
	return value, decode(val, value)
}

func decode(input interface{}, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			CustomHookFunc,
		),
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
