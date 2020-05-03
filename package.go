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
	"os"
	"path/filepath"
	"strings"
)

const defaultPath = "?.lua"
const pathListSeparator = ';'

func findLoader(l *lua.State, name string) {
	var msg string
	if l.Field(lua.UpValueIndex(1), "searchers"); !l.IsTable(3) {
		lua.Errorf(l, "'package.searchers' must be a table")
	}
	for i := 1; ; i++ {
		if l.RawGetInt(3, i); l.IsNil(-1) {
			l.Pop(1)
			l.PushString(msg)
			lua.Errorf(l, "module '%s' not found: %s", name, msg)
		}
		l.PushString(name)
		if l.Call(1, 2); l.IsFunction(-2) {
			return
		} else if l.IsString(-2) {
			msg += lua.CheckString(l, -2)
		}
		l.Pop(2)
	}
}

func findFile(l *lua.State, name, field, dirSep string) (string, error) {
	l.Field(lua.UpValueIndex(1), field)
	path, ok := l.ToString(-1)
	if !ok {
		lua.Errorf(l, "'package.%s' must be a string", field)
	}
	return searchPath(l, name, path, ".", dirSep)
}

func checkLoad(l *lua.State, loaded bool, fileName string) int {
	if loaded { // Module loaded successfully?
		l.PushString(fileName) // Second argument to module.
		return 2               // Return open function & file name.
	}
	m := lua.CheckString(l, 1)
	e := lua.CheckString(l, -1)
	lua.Errorf(l, "error loading module '%s' from file '%s':\n\t%s", m, fileName, e)
	panic("unreachable")
}

func searcherLua(l *lua.State) int {
	name := lua.CheckString(l, 1)
	filename, err := findFile(l, name, "path", string(filepath.Separator))
	if err != nil {
		return 1 // Module not found in this path.
	}
	return checkLoad(l, LoadFile(l, filename, "text") == nil, filename)
}

func searcherLibs(l *lua.State) int {
	name := lua.CheckString(l, 1)
	env, err := GetEnvironment(l)
	if err != nil {
		return 1
	}
	file, ok := env.LuaLibs[name]
	if !ok {
		return 1
	}

	return checkLoad(l, LoadLuaFile(l, file) == nil, file.Name)
}

func searcherPreload(l *lua.State) int {
	name := lua.CheckString(l, 1)
	l.Field(lua.RegistryIndex, "_PRELOAD")
	l.Field(-1, name)
	if l.IsNil(-1) {
		l.PushString(fmt.Sprintf("\n\tno field package.preload['%s']", name))
	}
	return 1
}

func createSearchersTable(l *lua.State) {
	searchers := []lua.Function{searcherPreload, searcherLibs, searcherLua}
	l.CreateTable(len(searchers), 0)
	for i, s := range searchers {
		l.PushValue(-2)
		l.PushGoClosure(s, 1)
		l.RawSetInt(-2, i+1)
	}
}

func searchPath(l *lua.State, name, path, sep, dirSep string) (string, error) {
	var msg string
	if sep != "" {
		name = strings.Replace(name, sep, dirSep, -1) // Replace sep by dirSep.
	}
	path = strings.Replace(path, string(pathListSeparator), string(filepath.ListSeparator), -1)
	for _, template := range filepath.SplitList(path) {
		if template != "" {
			filename := strings.Replace(template, "?", name, -1)
			msg = fmt.Sprintf("%s\n\tno file '%s'", msg, filename)
			env, err := GetEnvironment(l)
			if err != nil {
				return "", err
			}
			reader, err := env.Fs.GetReader(filename)
			if err != nil {
				return "", errors.New(msg)
			}
			err = reader.Close()
			if err != nil {
				return "", err
			}
			return filename, nil
		}
	}
	return "", errors.New(msg)
}

func noEnv(l *lua.State) bool {
	l.Field(lua.RegistryIndex, "LUA_NOENV")
	b := l.ToBoolean(-1)
	l.Pop(1)
	return b
}

func setPath(l *lua.State, field, env, def string) {
	if path := os.Getenv(env); path == "" || noEnv(l) {
		l.PushString(def)
	} else {
		o := fmt.Sprintf("%c%c", pathListSeparator, pathListSeparator)
		n := fmt.Sprintf("%c%s%c", pathListSeparator, def, pathListSeparator)
		path = strings.Replace(path, o, n, -1)
		l.PushString(path)
	}
	l.SetField(-2, field)
}

var packageLibrary = []lua.RegistryFunction{
	{"loadlib", func(l *lua.State) int {
		_ = lua.CheckString(l, 1) // path
		_ = lua.CheckString(l, 2) // init
		l.PushNil()
		l.PushString("dynamic libraries not enabled; check your Lua installation")
		l.PushString("absent")
		return 3 // Return nil, error message, and where.
	}},
	{"searchpath", func(l *lua.State) int {
		name := lua.CheckString(l, 1)
		path := lua.CheckString(l, 2)
		sep := lua.OptString(l, 3, ".")
		dirSep := lua.OptString(l, 4, string(filepath.Separator))
		f, err := searchPath(l, name, path, sep, dirSep)
		if err != nil {
			l.PushNil()
			l.PushString(err.Error())
			return 2
		}
		l.PushString(f)
		return 1
	}},
}

// PackageOpen opens the package library. Usually passed to Require.
func PackageOpen(l *lua.State) int {
	lua.NewLibrary(l, packageLibrary)
	createSearchersTable(l)
	l.SetField(-2, "searchers")
	setPath(l, "path", "LUA_PATH", defaultPath)
	l.PushString(fmt.Sprintf("%c\n%c\n?\n!\n-\n", filepath.Separator, pathListSeparator))
	l.SetField(-2, "config")
	lua.SubTable(l, lua.RegistryIndex, "_LOADED")
	l.SetField(-2, "loaded")
	lua.SubTable(l, lua.RegistryIndex, "_PRELOAD")
	l.SetField(-2, "preload")
	l.PushGlobalTable()
	l.PushValue(-2)
	lua.SetFunctions(l, []lua.RegistryFunction{{"require", func(l *lua.State) int {
		name := lua.CheckString(l, 1)
		l.SetTop(1)
		l.Field(lua.RegistryIndex, "_LOADED")
		l.Field(2, name)
		if l.ToBoolean(-1) {
			return 1
		}
		l.Pop(1)
		findLoader(l, name)
		l.PushString(name)
		l.Insert(-2)
		l.Call(2, 1)
		if !l.IsNil(-1) {
			l.SetField(2, name)
		}
		l.Field(2, name)
		if l.IsNil(-1) {
			l.PushBoolean(true)
			l.PushValue(-1)
			l.SetField(2, name)
		}
		return 1
	}}}, 1)
	l.Pop(1)
	return 1
}
