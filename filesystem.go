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
