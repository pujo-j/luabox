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

import "github.com/Shopify/go-lua"

func boxPrint(l *lua.State) int {
	env, err := GetEnvironment(l)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	n := l.Top()
	l.Global("tostring")
	for i := 1; i <= n; i++ {
		l.PushValue(-1) // function to be called
		l.PushValue(i)  // value to box_print
		l.Call(1, 1)
		s, ok := l.ToString(-1)
		if !ok {
			lua.Errorf(l, "'tostring' must return a string to 'box_print'")
			panic("unreachable")
		}
		if i > 1 {
			_, _ = env.Output.Write([]byte("\t"))
		}
		_, _ = env.Output.Write([]byte(s))
		l.Pop(1) // pop result
	}
	_, _ = env.Output.Write([]byte("\n"))
	return 0
}
