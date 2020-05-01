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
