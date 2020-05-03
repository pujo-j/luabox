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
	"bufio"
	"fmt"
	"github.com/Shopify/go-lua"
	"io"
	"strings"
)

func LoadLuaFile(l *lua.State, f LuaFile) error {
	var fileName = f.Name
	fileNameIndex := l.Top() + 1
	fileError := func(what string) error {
		fileName, _ := l.ToString(fileNameIndex)
		l.PushFString("cannot %s %s", what, fileName[1:])
		l.Remove(fileNameIndex)
		return lua.FileError
	}

	l.PushString("@" + fileName)
	r := bufio.NewReader(strings.NewReader(f.Code))
	if skipped, err := skipComment(r); err != nil {
		l.SetTop(fileNameIndex)
		return fileError("read")
	} else if skipped {
		r = bufio.NewReader(io.MultiReader(strings.NewReader("\n"), r))
	}
	s, _ := l.ToString(-1)
	err := l.Load(r, s, "text")
	switch err {
	case nil, lua.SyntaxError, lua.MemoryError: // do nothing
	default:
		l.SetTop(fileNameIndex)
		return fileError("read")
	}
	l.Remove(fileNameIndex)
	return err
}

func LoadFile(l *lua.State, fileName, mode string) error {
	if !("text" == mode) {
		return fmt.Errorf("invalid file mode %s", mode)
	}
	env, err := GetEnvironment(l)
	if err != nil {
		return err
	}
	fileNameIndex := l.Top() + 1
	fileError := func(what string) error {
		fileName, _ := l.ToString(fileNameIndex)
		l.PushFString("cannot %s %s", what, fileName[1:])
		l.Remove(fileNameIndex)
		return lua.FileError
	}
	var f io.Reader
	if fileName == "" {
		l.PushString("=stdin")
		f = env.Input
	} else {
		l.PushString("@" + fileName)
		var err error
		if f, err = env.Fs.GetReader(fileName); err != nil {
			return fileError("open")
		}
	}
	r := bufio.NewReader(f)
	if skipped, err := skipComment(r); err != nil {
		l.SetTop(fileNameIndex)
		return fileError("read")
	} else if skipped {
		r = bufio.NewReader(io.MultiReader(strings.NewReader("\n"), r))
	}
	s, _ := l.ToString(-1)
	err = l.Load(r, s, mode)
	if f != env.Input {
		_ = f.(io.ReadCloser).Close()
	}
	switch err {
	case nil, lua.SyntaxError, lua.MemoryError: // do nothing
	default:
		l.SetTop(fileNameIndex)
		return fileError("read")
	}
	l.Remove(fileNameIndex)
	return err
}

func skipComment(r *bufio.Reader) (bool, error) {
	bom := "\xEF\xBB\xBF"
	if ba, err := r.Peek(len(bom)); err != nil && err != io.EOF {
		return false, err
	} else if string(ba) == bom {
		_, _ = r.Read(ba)
	}
	if c, _, err := r.ReadRune(); err != nil {
		if err == io.EOF {
			err = nil
		}
		return false, err
	} else if c == '#' {
		_, err = r.ReadBytes('\n')
		if err == io.EOF {
			err = nil
		}
		return true, err
	}
	return false, r.UnreadRune()
}

func loadHelper(l *lua.State, s error, e int) int {
	if s == nil {
		if e != 0 {
			l.PushValue(e)
			if _, ok := lua.SetUpValue(l, -2, 1); !ok {
				l.Pop(1)
			}
		}
		return 1
	}
	l.PushNil()
	l.Insert(-2)
	return 2
}
