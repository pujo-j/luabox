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
	zap *zap.SugaredLogger
	cfg *zap.Config
}

func NewZapLog(cfg zap.Config) (*ZapLog, error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLog{logger.Sugar(), &cfg}, err
}

func (z *ZapLog) GetLevel() luabox.Level {
	switch z.cfg.Level.Level() {
	case zap.DebugLevel:
		return luabox.DebugLevel
	case zap.InfoLevel:
		return luabox.InfoLevel
	case zap.WarnLevel:
		return luabox.WarnLevel
	case zap.ErrorLevel:
		return luabox.ErrorLevel
	case zap.DPanicLevel:
		return luabox.FatalLevel
	case zap.PanicLevel:
		return luabox.FatalLevel
	case zap.FatalLevel:
		return luabox.FatalLevel
	default:
		return luabox.FatalLevel
	}
}

func (z *ZapLog) WithFields(context map[string]interface{}) luabox.Log {
	args := mapToArgs(context)
	if args != nil {
		return &ZapLog{z.zap.With(args...), z.cfg}
	} else {
		return &ZapLog{z.zap, z.cfg}
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
		z.zap.With(args...).Debug(msg)
	} else {
		z.zap.Debug(msg)
	}
}

func (z *ZapLog) Info(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.zap.With(args...).Info(msg)
	} else {
		z.zap.Info(msg)
	}
}

func (z *ZapLog) Warn(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.zap.With(args...).Warn(msg)
	} else {
		z.zap.Warn(msg)
	}
}

func (z *ZapLog) Error(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.zap.With(args...).Error(msg)
	} else {
		z.zap.Error(msg)
	}
}

func (z *ZapLog) Fatal(msg string, context map[string]interface{}) {
	args := mapToArgs(context)
	if args != nil {
		z.zap.With(args...).Fatalf(msg)
	} else {
		z.zap.Fatalf(msg)
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
