package memory

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func home() string {
	home := os.Getenv("PAC8_HOME")
	if home == "" {
		home = "."
	}
	return home
}

func LoadROM(e *[]error, path string, checksum string) *ROM {
	filename := filepath.Join(home(), path)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		*e = append(*e, err)
		return nil
	}
	rom := NewROM(data)
	romChecksum := rom.Checksum()
	if checksum != romChecksum {
		*e = append(*e, fmt.Errorf("invalid checksum for file: %s\nexpected: %v\nreceived: %v", filename, romChecksum, checksum))
	}
	return rom
}
