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
	"io"
	"time"
)

type FileInfo struct {
	Name         string
	SelfUrl      string
	IsDir        bool
	LastModified time.Time
	Size         uint64
	ETag         string
}

type FsError string

func (e FsError) Error() string {
	return string(e)
}

const EReadonly = FsError("readonly filesystem")

type Filesystem interface {
	GetReader(file string) (io.ReadCloser, error)
	GetWriter(file string) (io.WriteCloser, error)
	List(path string) ([]FileInfo, error)
	Delete(path string) error
}
