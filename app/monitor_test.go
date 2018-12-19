package app

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/blackchip-org/pac8/component/memory"
	. "github.com/blackchip-org/pac8/expect"
	"github.com/blackchip-org/pac8/machine"
)

type fixture struct {
	out    bytes.Buffer
	mon    *Monitor
	cursor *memory.Cursor
}

func newTestMonitor() *fixture {
	f := &fixture{}
	sys := newFixtureCab(nil)
	mach := machine.New(sys)
	f.mon = NewMonitor(mach)
	f.cursor = memory.NewCursor(mach.Mem)
	f.mon.out.SetOutput(&f.out)
	return f
}

func testMonitorInput(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func testMonitorRun(mon *Monitor) {
	go mon.Run()
	mon.mach.Run()
}

func TestBreakpointOn(t *testing.T) {
	f := newTestMonitor()
	f.cursor.PutN(0x01, 0x01, 0x01)
	f.mon.in = testMonitorInput("b 0x02 on \n g")
	testMonitorRun(f.mon)

	WithFormat(t, "%04x").Expect(f.mon.cpu.PC()).ToBe(0x0002)
}

func TestBreakpointOff(t *testing.T) {
	f := newTestMonitor()
	f.cursor.PutN(0x01, 0x01, 0x01)
	f.mon.in = testMonitorInput("b 0x02 on \n b 0x02 off \n g")
	testMonitorRun(f.mon)

	WithFormat(t, "%04x").Expect(f.mon.cpu.PC()).NotToBe(0x0002)
}

func TestDisassembleFirstLine(t *testing.T) {
	f := newTestMonitor()
	f.cursor.PutN(0x20, 0xcd, 0xab)
	f.mon.in = testMonitorInput("d \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(f.out.String(), "\n")
	fmt.Println(Dump(f.mon.mem, 0, 0x0f, AsciiDecoder))
	With(t).Expect(lines[0]).ToBe(
		"$0000:  20 cd ab  i20 $abcd",
	)
}

func TestDisassembleLastLine(t *testing.T) {
	f := newTestMonitor()
	f.cursor.Pos = 0x3f
	f.cursor.PutN(0x20, 0xcd, 0xab)
	f.mon.in = testMonitorInput("d \n q")
	testMonitorRun(f.mon)
	out := strings.TrimSpace(f.out.String())
	lines := strings.Split(out, "\n")
	fmt.Println(Dump(f.mon.mem, 0, 0x0f, AsciiDecoder))
	With(t).Expect(lines[len(lines)-1]).ToBe(
		"$003f:  20 cd ab  i20 $abcd",
	)
}

func TestDisassembleAt(t *testing.T) {
	f := newTestMonitor()
	f.cursor.Pos = 0x100
	f.cursor.PutN(0x20, 0xcd, 0xab)
	f.mon.in = testMonitorInput("d 0100 \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(f.out.String(), "\n")
	With(t).Expect(lines[0]).ToBe(
		"$0100:  20 cd ab  i20 $abcd",
	)
}

func TestDisassembleRange(t *testing.T) {
	f := newTestMonitor()
	f.cursor.Pos = 0x0112
	f.cursor.PutN(0x20, 0xcd, 0xab)
	f.mon.in = testMonitorInput("d 0100 0112 \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(strings.TrimSpace(f.out.String()), "\n")
	With(t).Expect(lines[len(lines)-1]).ToBe(
		"$0112:  20 cd ab  i20 $abcd",
	)
}

func TestGoContinued(t *testing.T) {
	// FIXME: There is a race condition here when restarting
	t.SkipNow()
	f := newTestMonitor()
	f.cursor.PutN(
		0x20, 0xcd, 0xab,
		0x21, 0x34, 0x12,
	)
	f.mon.in = testMonitorInput("b 0003 on \n g")
	testMonitorRun(f.mon)
	f.mon.in = testMonitorInput("b 0003 off \n b 0006 on \n g")
	testMonitorRun(f.mon)
	WithFormat(t, "%04x").Expect(f.mon.cpu.PC()).ToBe(0x0006)
}

func TestGoAddress(t *testing.T) {
	f := newTestMonitor()
	f.cursor.Pos = 0x100
	f.cursor.PutN(
		0x20, 0xcd, 0xab,
	)
	f.mon.in = testMonitorInput("b 103 on \n g 0100")
	testMonitorRun(f.mon)
	WithFormat(t, "%04x").Expect(f.mon.cpu.PC()).ToBe(0x103)
}

func TestMemoryFirstLine(t *testing.T) {
	f := newTestMonitor()
	f.mon.in = testMonitorInput("m \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(f.out.String(), "\n")
	With(t).Expect(lines[0]).ToBe(
		"$0000 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
	)
}

func TestMemoryLastLine(t *testing.T) {
	f := newTestMonitor()
	f.mon.in = testMonitorInput("m \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(strings.TrimSpace(f.out.String()), "\n")
	With(t).Expect(lines[len(lines)-1]).ToBe(
		"$00f0 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
	)
}

func TestMemoryPage(t *testing.T) {
	f := newTestMonitor()
	f.mon.in = testMonitorInput("m 0100 \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(strings.TrimSpace(f.out.String()), "\n")
	want := "$01f0 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryNextPage(t *testing.T) {
	f := newTestMonitor()
	f.mon.in = testMonitorInput("m 0100 \n m \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(strings.TrimSpace(f.out.String()), "\n")
	want := "$02f0 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestMemoryRange(t *testing.T) {
	f := newTestMonitor()
	f.mon.in = testMonitorInput("m 0100 018f \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(strings.TrimSpace(f.out.String()), "\n")
	want := "$0180 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................"
	have := lines[len(lines)-1]
	if want != have {
		t.Errorf("\n want: %v \n have: %v \n", want, have)
	}
}

func TestPoke(t *testing.T) {
	f := newTestMonitor()
	f.mon.in = testMonitorInput("p 0900 ab \n q")
	testMonitorRun(f.mon)
	WithFormat(t, "%02x").Expect(f.mon.mem.Load(0x900)).ToBe(0xab)
}

func TestPokeN(t *testing.T) {
	f := newTestMonitor()
	f.mon.in = testMonitorInput("p 0900 ab cd ef 12 34 \n q")
	testMonitorRun(f.mon)
	WithFormat(t, "%02x").Expect(f.mon.mem.Load(0x904)).ToBe(0x34)
}

func TestPeek(t *testing.T) {
	f := newTestMonitor()
	f.mon.mem.Store(0x0900, 0xab)
	f.mon.in = testMonitorInput("p 0900 \n q")
	testMonitorRun(f.mon)
	lines := strings.Split(f.out.String(), "\n")
	With(t).Expect(lines[0]).ToBe("$ab +171")
}

func TestTrace(t *testing.T) {
	f := newTestMonitor()
	f.cursor.PutN(
		0x20, 0x34, 0x12,
		0x10, 0x56,
	)
	f.mon.mach.Breakpoints[0x0005] = struct{}{}
	f.mon.in = testMonitorInput("t \n g")
	testMonitorRun(f.mon)
	lines := strings.Split(strings.TrimSpace(f.out.String()), "\n")
	fmt.Println(lines[1])
	With(t).Expect(lines[0]).ToBe(
		"$0000:  20 34 12  i20 $1234",
	)
	With(t).Expect(lines[1]).ToBe(
		"$0003:  10 56     i10 $56",
	)
}

func TestTraceOff(t *testing.T) {
	f := newTestMonitor()
	f.cursor.PutN(
		0x20, 0x34, 0x12,
		0x10, 0x56,
	)
	f.mon.mach.Breakpoints[0x0005] = struct{}{}
	f.mon.in = testMonitorInput("t \n q")
	testMonitorRun(f.mon)
	f.mon.in = testMonitorInput("t \n g")
	testMonitorRun(f.mon)
	lines := strings.Split(strings.TrimSpace(f.out.String()), "\n")
	With(t).Expect(lines[0]).ToBe("[break]")
}

func TestDump(t *testing.T) {
	var dumpTests = []struct {
		name     string
		start    int
		data     func() []int
		showFrom int
		showTo   int
		want     string
	}{
		{
			"one line", 0x10,
			func() []int { return []int{} },
			0x10, 0x20, "" +
				"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
		}, {
			"two lines", 0x10,
			func() []int { return []int{} },
			0x10, 0x30, "" +
				"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................\n" +
				"$0020 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
		}, {
			"jagged top", 0x10,
			func() []int { return []int{} },
			0x14, 0x30, "" +
				"$0010             00 00 00 00  00 00 00 00 00 00 00 00     ............\n" +
				"$0020 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................",
		}, {
			"jagged bottom", 0x10,
			func() []int { return []int{} },
			0x10, 0x2b, "" +
				"$0010 00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00 ................\n" +
				"$0020 00 00 00 00 00 00 00 00  00 00 00 00             ............",
		},
		{
			"single value", 0x10,
			func() []int { return []int{0, 0x41} },
			0x11, 0x11, "" +
				"$0010    41                                             A",
		},
		{
			"$40-$5f", 0x10,
			func() []int {
				data := make([]int, 0)
				for i := 0x40; i < 0x60; i++ {
					data = append(data, i)
				}
				return data
			},
			0x10, 0x30, "" +
				"$0010 40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f @ABCDEFGHIJKLMNO\n" +
				"$0020 50 51 52 53 54 55 56 57  58 59 5a 5b 5c 5d 5e 5f PQRSTUVWXYZ[\\]^_",
		},
	}

	m := memory.NewRAM(0x100)
	for _, test := range dumpTests {
		t.Run(test.name, func(t *testing.T) {
			for i, value := range test.data() {
				m.Store(uint16(test.start+i), uint8(value))
			}
			have := Dump(m, uint16(test.showFrom), uint16(test.showTo),
				AsciiDecoder)
			have = strings.TrimSpace(have)
			With(t).Expect(have).ToBe(test.want)
		})
	}
}
