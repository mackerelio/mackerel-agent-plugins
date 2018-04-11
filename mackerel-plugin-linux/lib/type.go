package mplinux

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type blockDevice struct {
	path string
}

func (b *blockDevice) name() string {
	d := strings.Split(b.path, string(os.PathSeparator))
	return d[len(d)-1]
}

func (b *blockDevice) isVirtual() bool {
	if strings.Index(b.path, "/devices/virtual/") != -1 {
		return true
	}
	return false
}

func (b *blockDevice) isRemovable() bool {
	content, err := ioutil.ReadFile(filepath.Join(b.path, "removable"))
	if err != nil {
		return false
	}
	if len(content) > 0 && string(content[0]) == "1" {
		return true
	}
	return false
}

func (b *blockDevice) stat() ([]string, error) {
	var stat []string

	content, err := ioutil.ReadFile(filepath.Join(b.path, "stat"))
	if err != nil {
		return stat, err
	}

	return strings.Fields(string(content)), nil
}
