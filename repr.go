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
	"bytes"
	"encoding/json"
	"github.com/Shopify/go-lua"
	"gopkg.in/yaml.v3"
)

func Repr(l *lua.State) int {
	args, err := PullVarargs(l, 1)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	switch len(args) {
	case 0:
		l.PushString("")
		return 0
	case 1:
		data, err := yaml.Marshal(args[0])
		if err != nil {
			l.PushString(err.Error())
			l.Error()
			return 0
		}
		l.PushString(string(data))
		return 1
	default:
		b := bytes.NewBuffer([]byte{})
		encoder := yaml.NewEncoder(b)
		for p := range args {
			b.WriteString("---\n")
			err := encoder.Encode(p)
			if err != nil {
				l.PushString(err.Error())
				l.Error()
				return 0
			}
		}
		l.PushString(b.String())
		return 1
	}
}

func Parse(l *lua.State) int {
	s := lua.CheckString(l, 1)
	var res interface{}
	err := yaml.Unmarshal([]byte(s), &res)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	DeepPush(l, res)
	return 1
}

func ParseJson(l *lua.State) int {
	s := lua.CheckString(l, 1)
	var res interface{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	DeepPush(l, res)
	return 1
}

func ReprJson(l *lua.State) int {
	p, err := PullTable(l, 1)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	data, err := json.Marshal(p)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	l.PushString(string(data))
	return 1
}
