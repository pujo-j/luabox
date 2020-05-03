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
	"github.com/markbates/pkger"
	"io/ioutil"
	"os"
	"strings"
)

func ScanLibs(libDir string) (map[string]LuaFile, error) {
	libs := make(map[string]LuaFile)
	err := pkger.Walk(libDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".lua") {
			luaFile, err := pkger.Open(path)
			if err != nil {
				return err
			}
			luaBytes, err := ioutil.ReadAll(luaFile)
			if err != nil {
				return err
			}
			var name string
			if strings.Contains(path, ":") {
				p := strings.Split(path, ":")[1]
				name = strings.TrimPrefix(strings.TrimSuffix(p, ".lua"), libDir)
			} else {
				name = strings.TrimPrefix(strings.TrimSuffix(path, ".lua"), libDir)
			}
			f := LuaFile{
				Name: name,
				Code: string(luaBytes),
			}
			libs[name] = f
		}
		return nil
	})
	return libs, err
}

var BaseLibs map[string]LuaFile

func init() {
	l, err := ScanLibs("/lua/")
	if err != nil {
		panic(err)
	}
	BaseLibs = l
}
