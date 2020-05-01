package localenv

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pujo-j/luabox"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Fs struct {
	BaseDir string
}

func (f *Fs) getPath(file string) (string, error) {
	file = path.Join(f.BaseDir, file)
	file = path.Clean(file)
	if !strings.HasPrefix(file, f.BaseDir) {
		return "", fmt.Errorf("invalid file path: %s", file)
	}
	return file, nil
}

func (f *Fs) GetReader(filePath string) (io.ReadCloser, error) {
	filePath, err := f.getPath(filePath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *Fs) GetWriter(filePath string) (io.WriteCloser, error) {
	filePath, err := f.getPath(filePath)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filePath, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *Fs) List(filePath string) ([]luabox.FileInfo, error) {
	filePath, err := f.getPath(filePath)
	if err != nil {
		return nil, err
	}
	list, err := ioutil.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	res := make([]luabox.FileInfo, 0)
	for _, e := range list {
		e2 := luabox.FileInfo{}
		e2.ETag = hex.EncodeToString(sha256.New().Sum([]byte(e.ModTime().Format(time.RFC3339Nano))))
		e2.LastModified = e.ModTime()
		e2.Name = e.Name()
		e2.Size = uint64(e.Size())
		e2.SelfUrl = path.Join(filePath, e2.Name)
		e2.IsDir = e.IsDir()
		res = append(res, e2)
	}
	return res, nil
}

func (f *Fs) Delete(filePath string) error {
	filePath, err := f.getPath(filePath)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	return err
}
