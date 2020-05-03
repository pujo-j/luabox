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

package test

import (
	"github.com/Shopify/go-lua"
	"github.com/pujo-j/luabox/localenv"
	"log"
	"os"
	"path"
	"runtime/pprof"
	"testing"
	"time"
)

func TestBase(t *testing.T) {
	f, err := os.Create("test_base.prof")
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()
	start := time.Now()
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	wd = path.Join(wd, "test")
	env, err := localenv.NewEnv(path.Join(wd, "lua"), path.Join(wd, "init"), []string{})
	t.Logf("Env created in %s", time.Since(start).String())
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	L, err := env.Init()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	start2 := time.Now()
	err = lua.DoString(L, "require('test_base')")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("Execution: %s", time.Since(start2).String())
	t.Logf("Total: %s", time.Since(start).String())
}
