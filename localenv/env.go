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

package localenv

import (
	"context"
	"github.com/Shopify/go-lua"
	"github.com/pujo-j/luabox"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
)

type ZapLog struct {
	Zap *zap.SugaredLogger
}

func NewZapLog(cfg zap.Config) (*ZapLog, error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLog{logger.Sugar()}, err
}

func (z *ZapLog) WithFields(context map[string]interface{}) luabox.Log {
	args := mapToArgs(context)
	if args != nil {
		return &ZapLog{z.Zap.With(args...)}
	} else {
		return &ZapLog{z.Zap}
	}
}

func mapToArgs(context map[string]interface{}) []interface{} {
	if context == nil || len(context) == 0 {
		return nil
	} else {
		var args []interface{}
		for k, v := range context {
			args = append(args, k)
			args = append(args, v)
		}
		return args
	}
}
func (z *ZapLog) Debug(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.Zap.With(args...).Debug(msg)
	} else {
		z.Zap.Debug(msg)
	}
}

func (z *ZapLog) Info(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.Zap.With(args...).Info(msg)
	} else {
		z.Zap.Info(msg)
	}
}

func (z *ZapLog) Warn(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.Zap.With(args...).Warn(msg)
	} else {
		z.Zap.Warn(msg)
	}
}

func (z *ZapLog) Error(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.Zap.With(args...).Error(msg)
	} else {
		z.Zap.Error(msg)
	}
}

func (z *ZapLog) Fatal(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.Zap.With(args...).Fatalf(msg)
	} else {
		z.Zap.Fatalf(msg)
	}
}

func NewEnv(basePath string, initScripts string, args []string) (*luabox.Environment, error) {
	config := zap.NewDevelopmentConfig()
	config.DisableCaller = true
	log, err := NewZapLog(config)
	if err != nil {
		return nil, err
	}
	env := make(map[string]string)
	for _, k := range os.Environ() {
		env[k] = os.Getenv(k)
	}
	fs := luabox.VFS{}
	fs.BaseFs = &Fs{BaseDir: path.Clean(basePath)}
	fs.Prefixes = map[string]luabox.Filesystem{}
	matches, err := filepath.Glob(initScripts + "/*.lua")
	if err != nil {
		return nil, err
	}
	sort.Strings(matches)
	var preInitLua []luabox.LuaFile
	for _, initFile := range matches {
		luaBytes, err := ioutil.ReadFile(initFile)
		if err != nil {
			log.Error(err.Error(), nil)
		} else {
			lf := luabox.LuaFile{Name: path.Base(initFile), Code: string(luaBytes)}
			preInitLua = append(preInitLua, lf)
		}
	}
	res := luabox.Environment{
		Fs:         &fs,
		Context:    context.Background(),
		Args:       args,
		Env:        env,
		Input:      os.Stdin,
		Output:     os.Stdout,
		Log:        log,
		GoLibs:     []lua.RegistryFunction{},
		PreInitLua: preInitLua,
		LuaLibs:    luabox.BaseLibs,
	}
	return &res, nil
}
