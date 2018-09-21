// https://github.com/descarte1/fuse-emulator-fuse/tree/fuse-1-3-6/z80/tests

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var out bytes.Buffer
var whitespace = regexp.MustCompile(" +")

func main() {
	out.WriteString("// Code generated by cpu/z80/fuse/gen.go. DO NOT EDIT.\n\n")
	out.WriteString("package z80\n")
	out.WriteString("import \"github.com/blackchip-org/pac8/memory\"\n")

	out.WriteString(`
	type fuseTest struct {
		name    string
		af      uint16
		bc      uint16
		de      uint16
		hl      uint16
		af1     uint16
		bc1     uint16
		de1     uint16
		hl1     uint16
		ix      uint16
		iy      uint16
		sp      uint16
		pc      uint16
		i       uint8
		r       uint8
		iff1    int
		iff2    int
		im      int
		halt    int
		tstates int

		snapshots []memory.Snapshot
		portReads []memory.Snapshot
		portWrites []memory.Snapshot
	}
	`)

	out.WriteString("var fuseTests = []fuseTest{\n")
	loadTests()
	out.WriteString("}\n\n")

	out.WriteString("var fuseResults = map[string]fuseTest{\n")
	loadResults()
	out.WriteString("}\n")

	err := ioutil.WriteFile("fuse_test.go", out.Bytes(), 0644)
	if err != nil {
		fatal("unable to save file", err)
	}
}

func loadTests() {
	file, err := os.Open("fuse/tests.in")
	if err != nil {
		fatal("unable to open", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		} else {
			name := line
			scanner.Scan()
			parseTest(name, scanner)
		}
	}
}

func loadResults() {
	file, err := os.Open("fuse/tests.expected")
	if err != nil {
		fatal("unable to open", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		} else {
			name := line
			out.WriteString("\"" + name + "\": ")
			line = ""
			scanner.Scan()
			parseTest(name, scanner)
		}
	}
}

func parseTest(name string, scanner *bufio.Scanner) {
	t := make(map[string]string)
	t["name"] = name

	// Scan for events (on expected results)
	portReads := []string{}
	portWrites := []string{}

	for {
		line := scanner.Text()
		// If the line does not start with a space, there are
		// no more events
		if !strings.HasPrefix(line, " ") {
			break
		}
		line = whitespace.ReplaceAllString(line, " ")
		f := strings.Fields(line)
		if f[1] == "PR" {
			portReads = append(portReads, fmt.Sprintf(
				"memory.Snapshot{Address: 0x%v, Values: []uint8{0x%v}},\n",
				f[2][2:4], f[3]))
		} else if f[1] == "PW" {
			portWrites = append(portWrites, fmt.Sprintf(
				"memory.Snapshot{Address: 0x%v, Values: []uint8{0x%v}},\n",
				f[2][2:4], f[3]))
		}
		scanner.Scan()
	}

	f1 := strings.Fields(scanner.Text())
	t["af"] = f1[0]
	t["bc"] = f1[1]
	t["de"] = f1[2]
	t["hl"] = f1[3]
	t["af1"] = f1[4]
	t["bc1"] = f1[5]
	t["de1"] = f1[6]
	t["hl1"] = f1[7]
	t["ix"] = f1[8]
	t["iy"] = f1[9]
	t["sp"] = f1[10]
	t["pc"] = f1[11]

	scanner.Scan()
	text2 := whitespace.ReplaceAllString(scanner.Text(), " ")
	f2 := strings.Fields(text2)
	t["i"] = f2[0]
	t["r"] = f2[1]
	t["iff1"] = f2[2]
	t["iff2"] = f2[3]
	t["im"] = f2[4]
	t["halt"] = f2[5]
	t["tstates"] = f2[6]

	t["snapshots"] = parseSnapshots(scanner)
	t["portReads"] = strings.Join(portReads, "")
	t["portWrites"] = strings.Join(portWrites, "")

	testTemplate.Execute(&out, t)
}

func parseSnapshots(scanner *bufio.Scanner) string {
	var tests bytes.Buffer
	for {
		scanner.Scan()
		line := strings.Fields(scanner.Text())
		if len(line) == 0 || line[0] == "-1" {
			break
		}
		address := line[0]
		values := []string{}
		for _, value := range line[1 : len(line)-1] {
			values = append(values, "0x"+value)
		}
		tests.WriteString(
			fmt.Sprintf("memory.Snapshot{Address: 0x%v, Values: []uint8{%v}},\n",
				address, strings.Join(values, ",")))
	}
	return tests.String()
}

func fatal(message string, err error) {
	fmt.Printf("error: %v: %v\n", message, err)
	os.Exit(1)
}

var testTemplate = template.Must(template.New("").Parse(`fuseTest{
	name: "{{.name}}",
	af: 0x{{.af}},
	bc: 0x{{.bc}},
	de: 0x{{.de}},
	hl: 0x{{.hl}},
	af1: 0x{{.af1}},
	bc1: 0x{{.bc1}},
	de1: 0x{{.de1}},
	hl1: 0x{{.hl1}},
	ix: 0x{{.ix}},
	iy: 0x{{.iy}},
	sp: 0x{{.sp}},
	pc: 0x{{.pc}},
	i: 0x{{.i}},
	r: 0x{{.r}},
	iff1: {{.iff1}},
	iff2: {{.iff2}},
	im: {{.im}},
	halt: {{.halt}},
	tstates: {{.tstates}},
	snapshots: []memory.Snapshot{
		{{.snapshots}}
	},
	portReads: []memory.Snapshot{
		{{.portReads}}
	},
	portWrites: []memory.Snapshot{
		{{.portWrites}}
	},
},
`))
