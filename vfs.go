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
	"encoding/hex"
	"hash/fnv"
	"io"
	"strings"
	"time"
)

type VFS struct {
	BaseFs   Filesystem
	Prefixes map[string]Filesystem
}

func (vfs *VFS) getFs(file string) (string, Filesystem) {
	split := strings.Split(file, "/")
	if len(split) < 2 {
		return file, vfs.BaseFs
	} else {
		prefix := split[0]
		fs, ok := vfs.Prefixes[prefix]
		if ok {
			newPath := strings.Join(split[1:], "/")
			return newPath, fs
		} else {
			return file, vfs.BaseFs
		}
	}
}
func (vfs *VFS) GetReader(file string) (io.ReadCloser, error) {
	file2, fs := vfs.getFs(file)
	return fs.GetReader(file2)
}

func (vfs *VFS) GetWriter(file string) (io.WriteCloser, error) {
	file2, fs := vfs.getFs(file)
	return fs.GetWriter(file2)
}

var utc, _ = time.LoadLocation("UTC")
var baseTime = time.Date(1970, time.January, 1, 0, 0, 0, 0, utc)

func (vfs *VFS) List(file string) ([]FileInfo, error) {
	if file == "/" || file == "" {
		list, err := vfs.BaseFs.List("/")
		if err != nil {
			return nil, err
		}
		for p := range vfs.Prefixes {
			etag := fnv.New64a()
			_, err := etag.Write([]byte(p))
			if err != nil {
				panic(err)
			}
			direntry := FileInfo{
				IsDir:        true,
				Name:         p,
				SelfUrl:      "/" + p,
				Size:         0,
				ETag:         hex.EncodeToString(etag.Sum([]byte{})),
				LastModified: baseTime}
			list = append(list, direntry)
		}
		return list, nil
	}
	file2, fs := vfs.getFs(file)
	return fs.List(file2)
}

func (vfs *VFS) Delete(file string) error {
	file2, fs := vfs.getFs(file)
	return fs.Delete(file2)
}
