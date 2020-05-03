/*
 *    Copyright 2020 Josselin Pujo
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package luabox

import (
	"errors"
	"fmt"
	"github.com/Shopify/go-lua"
	"reflect"
)

func PullStringTable(l *lua.State, idx int) (map[string]string, error) {
	if !l.IsTable(idx) {
		return nil, fmt.Errorf("need a table at index %d, got %s", idx, lua.TypeNameOf(l, idx))
	}

	// Table at idx
	l.PushNil() // Add free slot for the value, +1

	table := make(map[string]string)
	// -1:nil, idx:table
	for l.Next(idx) {
		// -1:val, -2:key, idx:table
		key, ok := l.ToString(-2)
		if !ok {
			return nil, fmt.Errorf("key should be a string (%v)", l.ToValue(-2))
		}
		val, ok := l.ToString(-1)
		if !ok {
			return nil, fmt.Errorf("value for key '%s' should be a string (%v)", key, l.ToValue(-1))
		}
		table[key] = val
		l.Pop(1) // remove val from top, -1
		// -1:key, idx: table
	}

	return table, nil
}

func PullTable(l *lua.State, idx int) (interface{}, error) {
	if !l.IsTable(idx) {
		return nil, fmt.Errorf("need a table at index %d, got %s", idx, lua.TypeNameOf(l, idx))
	}

	return pullTableRec(l, idx)
}

func pullTableRec(l *lua.State, idx int) (interface{}, error) {
	if !l.CheckStack(2) {
		return nil, errors.New("pull table, stack exhausted")
	}

	idx = l.AbsIndex(idx)
	if isArray(l, idx) {
		return pullArrayRec(l, idx)
	}

	table := make(map[string]interface{})

	l.PushNil()
	for l.Next(idx) {
		// -1: value, -2: key, ..., idx: table
		key, ok := l.ToString(-2)
		if !ok {
			err := fmt.Errorf("key should be a string (%s)", lua.TypeNameOf(l, -2))
			l.Pop(2)
			return nil, err
		}

		value, err := toGoValue(l, -1)
		if err != nil {
			l.Pop(2)
			return nil, err
		}

		table[key] = value

		l.Pop(1)
	}

	return table, nil
}

func isArray(l *lua.State, idx int) bool {
	if !l.IsTable(idx) {
		return false
	}
	l.Length(idx)
	length, ok := l.ToInteger(-1)
	l.Pop(1)
	if ok && length > 0 {
		return true
	}
	return false
}

func pullArrayRec(l *lua.State, idx int) (interface{}, error) {
	table := make([]interface{}, lua.LengthEx(l, idx))

	l.PushNil()
	for l.Next(idx) {
		k, ok := l.ToInteger(-2)
		if !ok {
			l.Pop(2)
			return nil, fmt.Errorf("pull array: expected numeric index, got '%s'", l.TypeOf(-2))
		}

		v, err := toGoValue(l, -1)
		if err != nil {
			l.Pop(2)
			return nil, err
		}

		table[k-1] = v
		l.Pop(1)
	}

	return table, nil
}

func toGoValue(l *lua.State, idx int) (interface{}, error) {
	t := l.TypeOf(idx)
	switch t {
	case lua.TypeBoolean:
		return l.ToBoolean(idx), nil
	case lua.TypeString:
		return lua.CheckString(l, idx), nil
	case lua.TypeNumber:
		return lua.CheckInteger(l, idx), nil
	case lua.TypeTable:
		return pullTableRec(l, idx)
	default:
		err := fmt.Errorf("pull table, unsupported type %s", lua.TypeNameOf(l, idx))
		return nil, err
	}
}

func PullVarargs(l *lua.State, startIndex int) ([]interface{}, error) {
	top := l.Top()
	if top < startIndex {
		return []interface{}{}, nil
	}

	varargs := make([]interface{}, top-startIndex+1)
	for i := startIndex; i <= top; i++ {
		var value interface{}
		var err error
		switch l.TypeOf(i) {
		case lua.TypeNil:
			value = nil
		case lua.TypeBoolean:
			value = l.ToBoolean(i)
		case lua.TypeLightUserData:
			value = nil // not supported by go-lua
		case lua.TypeNumber:
			value = lua.CheckNumber(l, i)
		case lua.TypeString:
			value = lua.CheckString(l, i)
		case lua.TypeTable:
			value, err = PullTable(l, i)
			if err != nil {
				return nil, err
			}
		case lua.TypeFunction:
			value = l.ToGoFunction(i)
		case lua.TypeUserData:
			value = l.ToUserData(i)
		case lua.TypeThread:
			value = l.ToThread(i)
		}
		varargs[i-startIndex] = value
	}
	return varargs, nil
}

func MustPullVarargs(l *lua.State, startIndex int) []interface{} {
	varargs, err := PullVarargs(l, startIndex)
	if err != nil {
		lua.Errorf(l, err.Error())
		panic("unreachable")
	}
	return varargs
}

func DeepPush(l *lua.State, v interface{}) int {
	forwardOnType(l, v)
	return 1
}

func forwardOnType(l *lua.State, val interface{}) {

	switch val := val.(type) {
	case nil:
		l.PushNil()

	case bool:
		l.PushBoolean(val)

	case string:
		l.PushString(val)

	case uint8:
		l.PushNumber(float64(val))
	case uint16:
		l.PushNumber(float64(val))
	case uint32:
		l.PushNumber(float64(val))
	case uint64:
		l.PushNumber(float64(val))
	case uint:
		l.PushNumber(float64(val))

	case int8:
		l.PushNumber(float64(val))
	case int16:
		l.PushNumber(float64(val))
	case int32:
		l.PushNumber(float64(val))
	case int64:
		l.PushNumber(float64(val))
	case int:
		l.PushNumber(float64(val))

	case float32:
		l.PushNumber(float64(val))
	case float64:
		l.PushNumber(val)

	case complex64:
		forwardOnType(l, []float32{real(val), imag(val)})
	case complex128:
		forwardOnType(l, []float64{real(val), imag(val)})

	default:
		forwardOnReflect(l, val)
	}
}

func forwardOnReflect(l *lua.State, val interface{}) {

	switch v := reflect.ValueOf(val); v.Kind() {

	case reflect.Array, reflect.Slice:
		recurseOnFuncSlice(l, func(i int) interface{} { return v.Index(i).Interface() }, v.Len())

	case reflect.Map:
		l.CreateTable(0, v.Len())
		for _, key := range v.MapKeys() {
			mapKey := key.Interface()
			mapVal := v.MapIndex(key).Interface()
			forwardOnType(l, mapKey)
			forwardOnType(l, mapVal)
			l.RawSet(-3)
		}

	default:
		lua.Errorf(l, fmt.Sprintf("contains unsupported type: %T", val))
		panic("unreachable")
	}

}

// the hack of using a func(int)interface{} makes it that it is valid for any
// type of slice
func recurseOnFuncSlice(l *lua.State, input func(int) interface{}, n int) {
	l.CreateTable(n, 0)
	for i := 0; i < n; i++ {
		forwardOnType(l, input(i))
		l.RawSetInt(-2, i+1)
	}
}
