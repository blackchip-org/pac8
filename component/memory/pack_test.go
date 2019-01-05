package memory

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/blackchip-org/pac8/expect"
)

func resetReadFile() {
	readFile = ioutil.ReadFile
}

func TestPack(t *testing.T) {
	data := []byte{1, 42}
	checksum := fmt.Sprintf("%040x", sha1.Sum(data))
	readFile = func(filename string) ([]byte, error) {
		return []byte{1, 42}, nil
	}
	defer resetReadFile()
	p := NewPack().Add("group1", "/foo", checksum)
	set, err := p.Load("")
	With(t).Expect(err).ToBe(nil)
	With(t).Expect(set["group1"].Load(1)).ToBe(42)
}

func TestPackConcat(t *testing.T) {
	data1 := []byte{1, 42}
	data2 := []byte{1, 44}
	checksum1 := fmt.Sprintf("%040x", sha1.Sum(data1))
	checksum2 := fmt.Sprintf("%040x", sha1.Sum(data2))
	readFile = func(filename string) ([]byte, error) {
		switch filename {
		case "data1":
			return data1, nil
		case "data2":
			return data2, nil
		}
		return nil, fmt.Errorf("invalid file")
	}
	defer resetReadFile()
	p := NewPack().
		Add("group1", "data1", checksum1).
		Add("group1", "data2", checksum2)
	set, err := p.Load("")
	With(t).Expect(err).ToBe(nil)
	With(t).Expect(set["group1"].Load(1)).ToBe(42)
	With(t).Expect(set["group1"].Load(3)).ToBe(44)
}
