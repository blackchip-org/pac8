package memory

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var readFile = ioutil.ReadFile

type packEntry struct {
	path     string
	checksum string
}

type Pack struct {
	entries map[string][]packEntry
}

type Set map[string]Memory

func NewPack() *Pack {
	return &Pack{entries: make(map[string][]packEntry)}
}

func (p *Pack) Add(group string, path string, checksum string) *Pack {
	group = strings.TrimSpace(group)
	path = strings.TrimSpace(path)
	entry, ok := p.entries[group]
	if !ok {
		entry = make([]packEntry, 0, 0)
	}
	entry = append(entry, packEntry{path, checksum})
	p.entries[group] = entry
	return p
}

func (p *Pack) Load(dir string) (Set, error) {
	roms := make(map[string]Memory)
	e := make([]string, 0, 0)
	for group, entries := range p.entries {
		var rom bytes.Buffer
		for _, entry := range entries {
			path := filepath.Join(dir, entry.path)
			data, err := readFile(path)
			if err != nil {
				e = append(e, err.Error())
				continue
			}
			checksum := fmt.Sprintf("%040x", sha1.Sum(data))
			if checksum != entry.checksum {
				e = append(e, fmt.Sprintf("%v: invalid checksum", path))
				continue
			}
			rom.Write(data)
		}
		roms[group] = NewROM(rom.Bytes())
	}
	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, "\n"))
	}
	return roms, nil
}
