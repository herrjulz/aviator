package mingoak

import (
	"os"
	"time"
)

type File struct {
	content []byte
	name    string
	time    time.Time
}

func (f File) IsDir() bool {
	return false
}

func (d File) Name() string {
	return d.name
}

func (d File) Size() int64 {
	return int64(len(d.content))
}

func (d File) Mode() os.FileMode {
	return 777
}

func (d File) ModTime() time.Time {
	return d.time
}

func (d File) Sys() interface{} {
	return nil
}
