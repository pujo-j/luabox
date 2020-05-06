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
	"fmt"
	"github.com/Shopify/go-lua"
	"strconv"
)

var luaBoxSyscalls = []lua.RegistryFunction{
	{"log", func(l *lua.State) int {
		level := lua.CheckInteger(l, 1)
		message := lua.CheckString(l, 2)
		env, err := GetEnvironment(l)
		if err != nil {
			l.PushString(err.Error())
			l.Error()
			return 0
		}
		log := env.Log
		p, err := PullTable(l, 3)
		var params map[string]interface{}
		if err != nil {
			// No params
			params = map[string]interface{}{}
		} else {
			params = p.(map[string]interface{})
		}
		if f, ok := lua.Stack(l, 2); ok { // check function at level
			ar, _ := lua.Info(l, "Sl", f) // get info about it
			if ar.CurrentLine > 0 {       // is there info?
				params["lua_origin"] = fmt.Sprintf("%s:%d", ar.ShortSource, ar.CurrentLine)
			}
		}
		switch level {
		case 1:
			log.Debug(message, params)
			break
		case 2:
			log.Info(message, params)
			break
		case 3:
			log.Warn(message, params)
			break
		case 4:
			log.Error(message, params)
			break
		case 5:
			log.Fatal(message, params)
			break
		default:
			l.PushString("invalid log level " + strconv.Itoa(level))
			l.Error()
			return 0
		}
		return 0
	}},
	{"yamlRepr", Repr},
	{"yamlParse", Parse},
	{"jsonRepr", ReprJson},
	{"jsonParse", ParseJson},
	{"getEnv", EnvGetEnv},
	{"getArgs", EnvGetArgs},
}

func SyscallOpen(l *lua.State) int {
	lua.NewLibrary(l, luaBoxSyscalls)
	return 1
}
