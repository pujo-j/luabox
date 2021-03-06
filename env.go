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
	"context"
	"errors"
	"github.com/Shopify/go-lua"
	"io"
	"reflect"
)

type LuaFile struct {
	Name string
	Code string
}

const EnvKey = "LUABOX_ENV"

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type Environment struct {
	Fs         Filesystem
	Input      io.Reader
	Output     WriteSyncer
	Log        Log
	Env        map[string]string
	Args       []string
	Context    context.Context
	GoLibs     []lua.RegistryFunction
	LuaLibs    map[string]LuaFile
	PreInitLua []LuaFile
}

func (e *Environment) Init() (*lua.State, error) {
	if e.Env == nil {
		e.Env = make(map[string]string)
	}
	if e.Args == nil {
		e.Args = make([]string, 0)
	}
	l := lua.NewState()
	// Get _G on stack
	lua.Require(l, "_G", lua.BaseOpen, true)
	l.PushGoFunction(func(l *lua.State) int {
		f := lua.OptString(l, 1, "")
		if l.SetTop(1); LoadFile(l, f, "") != nil {
			l.Error()
			panic("unreachable")
		}
		continuation := func(l *lua.State) int { return l.Top() - 1 }
		l.CallWithContinuation(0, lua.MultipleReturns, 0, continuation)
		return continuation(l)
	})
	l.SetField(-2, "dofile")
	l.PushGoFunction(boxPrint)
	l.SetField(-2, "print")
	l.PushGoFunction(func(l *lua.State) int {
		f, m, e := lua.OptString(l, 1, ""), lua.OptString(l, 2, ""), 3
		if l.IsNone(e) {
			e = 0
		}
		return loadHelper(l, LoadFile(l, f, m), e)
	})
	l.SetField(-2, "loadfile")
	libs := []lua.RegistryFunction{
		{"package", PackageOpen},
		{"table", lua.TableOpen},
		{"string", lua.StringOpen},
		{"bit32", lua.Bit32Open},
		{"math", lua.MathOpen},
		{"luabox", SyscallOpen},
	}
	for _, lib := range libs {
		lua.Require(l, lib.Name, lib.Function, true)
		l.Pop(1)
	}
	// Expose go libraries
	if e.GoLibs != nil {
		for _, lib := range e.GoLibs {
			lua.Require(l, lib.Name, lib.Function, true)
			l.Pop(1)
		}
	}
	SetEnvironment(l, e)
	if e.PreInitLua != nil {
		for _, script := range e.PreInitLua {
			err := lua.LoadBuffer(l, script.Code, script.Name, "text")
			if err != nil {
				e.Log.Error("loading preinit", map[string]interface{}{"script": script.Name, "error": err.Error()})
			} else {
				err := l.ProtectedCall(0, 0, 0)
				if err != nil {
					e.Log.Error("loading preinit", map[string]interface{}{"script": script.Name, "error": err.Error()})
				}
			}
		}
	}
	// And finally remove _G from stack
	l.Pop(1)
	return l, nil
}

func SetEnvironment(l *lua.State, environment *Environment) {
	l.PushString(EnvKey)
	l.PushUserData(environment)
	l.SetTable(lua.RegistryIndex)
}

func GetEnvironment(l *lua.State) (*Environment, error) {
	l.PushString(EnvKey)
	l.Table(lua.RegistryIndex)
	ud := l.ToUserData(-1)
	l.Pop(1)
	res, ok := ud.(*Environment)
	if !ok {
		return nil, errors.New("invalid type:" + reflect.TypeOf(ud).Name())
	}
	return res, nil
}

func EnvGetEnv(l *lua.State) int {
	environment, err := GetEnvironment(l)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	DeepPush(l, environment.Env)
	return 1
}

func EnvGetArgs(l *lua.State) int {
	environment, err := GetEnvironment(l)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
		return 0
	}
	DeepPush(l, environment.Args)
	return 1
}
