package app

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/blackchip-org/pac8/memory"
)

const (
	Name = "Portable Arcade Cabinet"
	Slug = "pac8"
)

var (
	home     string
	userHome string
)

func init() {
	userHome = "."
	u, err := user.Current()
	if err != nil {
		log.Printf("unable to find home directory: %v", err)
	} else {
		userHome = u.HomeDir
	}
	flag.StringVar(&home, "home", "", "path to runtime data")
}

var (
	// ROM is the directory with ROM images
	ROM = "rom"
	// Ext is the directory with external data (test data, docs...)
	Ext = "ext"
	// Store is the directory with variable runtime data (high scores...)
	Store = "var"
)

func PathFor(kind string, path ...string) string {
	root := home
	if root == "" {
		envVar := fmt.Sprintf("%v_HOME", strings.ToUpper(Slug))
		root = os.Getenv(envVar)
	}
	if root == "" {
		root = filepath.Join(userHome, Slug)
	}
	return filepath.Join(root, kind, filepath.Join(path...))
}

func LoadROM(e *[]error, path string, checksum string) memory.Memory {
	filename := PathFor(ROM, path)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		*e = append(*e, err)
		return nil
	}
	rom := memory.NewROM(data)
	romChecksum := fmt.Sprintf("%04x", sha1.Sum(data))
	if checksum != romChecksum {
		*e = append(*e, fmt.Errorf("%v: invalid checksum", filename))
	}
	return rom
}

type Loader struct {
	e   []error
	dir string
}

func NewLoader(dir string) *Loader {
	return &Loader{
		e:   make([]error, 0, 0),
		dir: dir,
	}
}

func (l *Loader) Load(name string, checksum string) memory.Memory {
	path := filepath.Join(l.dir, name)
	return LoadROM(&l.e, path, checksum)
}

func (l *Loader) Error() error {
	if len(l.e) == 0 {
		return nil
	}
	var out bytes.Buffer
	for i, e := range l.e {
		out.WriteString(e.Error())
		if i != len(l.e)-1 {
			out.WriteString("\n")
		}
	}
	return errors.New(out.String())
}
