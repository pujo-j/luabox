package luabox

import (
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"io/ioutil"
	"strings"
)

func ScanLibs(libDir string, libDirName string) (map[string]LuaFile, error) {
	p := packr.New(libDirName, libDir)
	libs := make(map[string]LuaFile)
	err := p.Walk(func(s string, file packd.File) error {
		if strings.HasSuffix(s, ".lua") {
			luaBytes, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			name := strings.Split(s, ".")[0]
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
	l, err := ScanLibs("lua", "baselibs")
	if err != nil {
		panic(err)
	}
	BaseLibs = l
}
