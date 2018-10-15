package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	err := toArray("zexdoc", "./zexdoc.com", "./zexdoc_com_test.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func toArray(name string, inFile string, outFile string) error {
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	out.WriteString("package main\n")
	out.WriteString(fmt.Sprintf("var %v = []byte{\n", name))
	for i, b := range data {
		out.WriteString(fmt.Sprintf("0x%02x,", b))
		if i%8 == 0 {
			out.WriteString("\n")
		}
	}
	out.WriteString("}")
	err = ioutil.WriteFile(outFile, out.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}
