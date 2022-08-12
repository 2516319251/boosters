package file

import (
	"errors"
	"io"
	"os"

	"github.com/2516319251/boosters/config"
)

var _ config.Source = (*file)(nil)

type file struct {
	path string
}

func Load(path string) config.Source {
	return &file{path: path}
}

func (file *file) Load() (*config.KeyValue, error) {
	fi, err := os.Stat(file.path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, errors.New("can not load dir")
	}

	return file.loadFile(file.path)
}

func (file *file) loadFile(path string) (*config.KeyValue, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &config.KeyValue{
		Key:    info.Name(),
		Value:  data,
		Format: format(info.Name()),
	}, nil
}
