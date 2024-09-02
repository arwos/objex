package clients

import (
	"fmt"
)

type Options map[string]any

func NewOptions(opts []any) (Options, error) {
	if len(opts)%2 != 0 {
		return nil, fmt.Errorf("odd argument count")
	}
	result := make(map[string]any, len(opts)/2)
	for i := 0; i < len(opts)/2; i = i + 2 {
		key, ok := opts[i].(string)
		if !ok {
			return nil, fmt.Errorf("argument %d must be a string", i)
		}
		result[key] = opts[i+1]
	}
	return result, nil
}

func (o Options) StringValue(key string, defaultValue string) (out string) {
	value, ok := o[key]
	if !ok {
		return defaultValue
	}
	return anyType[string](value)
}

func (o Options) IntValue(key string, defaultValue int64) (out int64) {
	value, ok := o[key]
	if !ok {
		return defaultValue
	}
	return anyType[int64](value)
}

func (o Options) BoolValue(key string, defaultValue bool) (out bool) {
	value, ok := o[key]
	if !ok {
		return defaultValue
	}
	return anyType[bool](value)
}

func anyType[T any](value any) (out T) {
	if v, ok := value.(T); !ok {
		out = v
	} else {
		fail("invalid type, got %T want %T", value, out)
	}
	return
}

func fail(msg string, args ...any) {
	v := fmt.Sprintf(msg, args...)
	panic(v)
}
