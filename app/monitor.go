package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc"
	"github.com/blackchip-org/pac8/machine"
	"github.com/chzyer/readline"
)

const (
	CmdBreakpoint  = "b"
	CmdDisassemble = "d"
	CmdFill        = "f"
	CmdGo          = "g"
	CmdHalt        = "h"
	CmdHelp        = "?"
	CmdMemory      = "m"
	CmdNext        = "n"
	CmdPokePeek    = "p"
	CmdRegisters   = "r"
	CmdStep        = "s"
	CmdRestore     = "si"
	CmdSave        = "so"
	CmdTrace       = "t"
	CmdQuit        = "q"
	CmdQuitLong    = "quit"
)

const (
	memPageLen       = 0x100
	dasmPageLen      = 0x3f
	maxArgs          = 0x100
	SnapshotFileName = "snapshot"
)

type CharDecoder func(uint8) (rune, bool)

type Monitor struct {
	dasm         *proc.Disassembler
	mach         *machine.Mach
	cpu          proc.CPU
	mem          memory.Memory
	breakpoints  map[uint16]struct{}
	in           io.ReadCloser
	out          *log.Logger
	rl           *readline.Instance
	lastCmd      string
	memPtr       uint16
	dasmPtr      uint16
	selectedCore int
}

func NewMonitor(mach *machine.Mach) *Monitor {
	m := &Monitor{
		mach:        mach,
		cpu:         mach.Cores[0].CPU,
		mem:         mach.Cores[0].Mem,
		breakpoints: mach.Cores[0].Breakpoints,
		in:          readline.NewCancelableStdin(os.Stdin),
		out:         log.New(os.Stdout, "", 0),
		dasm:        mach.Cores[0].CPU.Info().NewDisassembler(mach.Cores[0].Mem),
	}
	mach.EventCallback = m.handleEvent
	return m
}

func (m *Monitor) Run() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      m.getPrompt(),
		HistoryFile: filepath.Join(usr.HomeDir, ".pac8-history"), // FIXME: Move this to runtime directory

		Stdin: m.in,
	})
	if err != nil {
		return err
	}
	m.rl = rl
	for {
		line, err := rl.Readline()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		m.parse(line)
	}
}

func (m *Monitor) Close() {
	m.in.Close()
}

func (m *Monitor) getPrompt() string {
	c := ""
	if len(m.mach.Cores) > 1 {
		c = fmt.Sprintf(":%v", m.selectedCore)
	}
	return fmt.Sprintf("monitor%v> ", c)
}

func (m *Monitor) parse(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		if m.lastCmd != CmdStep && m.lastCmd != CmdGo {
			return
		}
		line = m.lastCmd
	}
	fields := strings.Split(line, " ")

	if len(fields) == 0 {
		return
	}

	cmd := fields[0]
	args := fields[1:]
	var err error
	switch cmd {
	case CmdBreakpoint:
		err = m.breakpoint(args)
	case CmdDisassemble:
		err = m.disassemble(args)
	case CmdFill:
		err = m.fill(args)
	case CmdGo:
		err = m.goCmd(args)
	case CmdHalt:
		err = m.halt(args)
	case CmdHelp:
		err = m.help(args)
	case CmdMemory:
		err = m.memory(args, m.mach.CharDecoder)
	case CmdNext:
		err = m.next(args)
	case CmdPokePeek:
		err = m.pokePeek(args)
	case CmdRegisters:
		err = m.registers(args)
	case CmdRestore:
		err = m.restore(args)
	case CmdSave:
		err = m.save(args)
	case CmdStep:
		err = m.step(args)
	case CmdTrace:
		err = m.trace(args)
	case CmdQuit, CmdQuitLong:
		m.rl.Close()
		m.mach.Send(machine.QuitCmd)
		runtime.Goexit()
	default:
		err = fmt.Errorf("unknown command: %v", cmd)
	}

	if err != nil {
		m.out.Println(err)
	} else {
		m.lastCmd = cmd
	}
}

func (m *Monitor) breakpoint(args []string) error {
	if err := checkLen(args, 1, 2); err != nil {
		return err
	}
	address, err := parseAddress(args[0])
	if err != nil {
		return err
	}
	if len(args) == 1 {
		if _, exists := m.breakpoints[address]; exists {
			m.out.Println("breakpoint on")
		} else {
			m.out.Println("breakpoint off")
		}
		return nil
	}
	switch args[1] {
	case "on":
		m.breakpoints[address] = struct{}{}
	case "off":
		delete(m.breakpoints, address)
	default:
		return fmt.Errorf("invalid: %v", args[1])
	}
	return nil
}

func (m *Monitor) disassemble(args []string) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}
	addrStart := m.cpu.PC()
	if len(args) == 0 {
		if strings.HasPrefix(m.lastCmd, CmdDisassemble) {
			addrStart = m.dasmPtr
		}
	}
	if len(args) > 0 {
		addr, err := parseAddress(args[0])
		if err != nil {
			return err
		}
		addrStart = addr
	}
	addrEnd := addrStart + uint16(dasmPageLen)
	if len(args) > 1 {
		addr, err := parseAddress(args[1])
		if err != nil {
			return err
		}
		addrEnd = addr
	}
	m.dasm.SetPC(addrStart)
	for m.dasm.PC() <= addrEnd {
		m.out.Println(m.dasm.Next())
	}
	m.dasmPtr = m.dasm.PC()
	return nil
}

func (m *Monitor) fill(args []string) error {
	if err := checkLen(args, 3, 3); err != nil {
		return err
	}
	start, err := parseAddress(args[0])
	if err != nil {
		return err
	}
	end, err := parseAddress(args[1])
	if err != nil {
		return err
	}
	value, err := parseValue(args[2])
	if err != nil {
		return err
	}
	for addr := start; addr <= end; addr++ {
		m.mem.Store(addr, value)
	}
	return nil
}

func (m *Monitor) goCmd(args []string) error {
	if err := checkLen(args, 0, 1); err != nil {
		return err
	}
	if len(args) > 0 {
		address, err := parseAddress(args[0])
		if err != nil {
			return err
		}
		m.cpu.SetPC(address)
	}
	go m.mach.Send(machine.StartCmd)
	return nil
}

func (m *Monitor) halt(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.mach.Send(machine.StopCmd)
	return nil
}

func (m *Monitor) help(args []string) error {
	if err := checkLen(args, 0, 1); err != nil {
		return err
	}
	if len(args) == 0 {
		m.out.Println(helpList)
		return nil
	}
	cmd := args[0]
	text, ok := helpCmds[cmd]
	if !ok {
		return errors.New("no such command")
	}
	m.out.Println(text)
	return nil
}

func (m *Monitor) memory(args []string, decoder CharDecoder) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}
	addrStart := m.cpu.PC()
	if len(args) == 0 {
		if m.lastCmd == CmdMemory {
			addrStart = m.memPtr
		}
	}
	if len(args) > 0 {
		addr, err := parseAddress(args[0])
		if err != nil {
			return err
		}
		addrStart = addr
	}
	addrEnd := addrStart + uint16(memPageLen)
	if len(args) > 1 {
		addr, err := parseAddress(args[1])
		if err != nil {
			return err
		}
		addrEnd = addr
	}
	m.out.Println(Dump(m.mem, addrStart, addrEnd, decoder))
	m.memPtr = addrEnd
	return nil
}

func (m *Monitor) next(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.dasm.SetPC(m.cpu.PC())
	m.out.Println(m.dasm.Next())
	return nil
}

func (m *Monitor) pokePeek(args []string) error {
	if err := checkLen(args, 1, maxArgs); err != nil {
		return err
	}
	address, err := parseAddress(args[0])
	if err != nil {
		return err
	}
	// peek
	if len(args) == 1 {
		m.out.Println(formatValue(m.mem.Load(address)))
		return nil
	}
	// poke
	values := []uint8{}
	for _, str := range args[1:] {
		v, err := parseValue(str)
		if err != nil {
			return err
		}
		values = append(values, v)
	}
	for offset, v := range values {
		m.mem.Store(address+uint16(offset), v)
	}
	return nil
}

func (m *Monitor) registers(args []string) error {
	if err := checkLen(args, 0, 2); err != nil {
		return err
	}

	// Print all registers
	if len(args) == 0 {
		m.out.Printf("[%v]\n", m.mach.Status)
		m.out.Println(m.cpu.String())
		return nil
	}

	name := strings.ToUpper(args[0])
	reg, ok := m.cpu.Info().Registers[name]
	if !ok {
		return errors.New("no such register")
	}

	// Get value of register
	if len(args) == 1 {
		switch get := reg.Get.(type) {
		case func() uint8:
			m.out.Println(formatValue(get()))
		case func() uint16:
			m.out.Println(formatValue16(get()))
		default:
			panic("unexpected type")
		}
		return nil
	}

	// Set value of register
	switch put := reg.Put.(type) {
	case func(uint8):
		v, err := parseValue(args[1])
		if err != nil {
			return nil
		}
		put(v)
	case func(uint16):
		v, err := parseValue16(args[1])
		if err != nil {
			return nil
		}
		put(v)
	}
	return nil
}

func (m *Monitor) restore(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	fileName := PathFor(Store, m.mach.System.Spec().Name, SnapshotFileName)
	m.mach.Send(machine.RestoreCmd, fileName)
	return nil
}

func (m *Monitor) save(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	fileName := PathFor(Store, m.mach.System.Spec().Name, SnapshotFileName)
	m.mach.Send(machine.SaveCmd, fileName)
	return nil
}

func (m *Monitor) step(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.cpu.Next()
	m.dasm.SetPC(m.cpu.PC())
	m.out.Println(m.dasm.Next())
	return nil
}

func (m *Monitor) trace(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.mach.Send(machine.TraceCmd)
	return nil
}

func (m *Monitor) handleEvent(evt machine.EventType, arg interface{}) {
	switch evt {
	case machine.StatusEvent:
		s := arg.(machine.Status)
		if s == machine.Break {
			fmt.Println()
			m.registers([]string{})
			m.rl.Refresh()
		}
	case machine.TraceEvent:
		msg := arg.(string)
		m.out.Println(msg)
	case machine.ErrorEvent:
		msg := arg.(string)
		m.out.Println(msg)
	default:
		log.Panicf("unknown arg: %v", arg)
	}
}

func checkLen(args []string, min int, max int) error {
	if len(args) < min {
		return errors.New("not enough arguments")
	}
	if len(args) > max {
		return errors.New("too many arguments")
	}
	return nil
}

func parseUint(str string, bitSize int) (uint64, error) {
	base := 16
	switch {
	case strings.HasPrefix(str, "$"):
		str = str[1:]
	case strings.HasPrefix(str, "0x"):
		str = str[2:]
	case strings.HasPrefix(str, "+"):
		str = str[1:]
		base = 10
	}
	return strconv.ParseUint(str, base, bitSize)
}

func parseAddress(str string) (uint16, error) {
	value, err := parseUint(str, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid address: %v", str)
	}
	return uint16(value), nil
}

func parseValue(str string) (uint8, error) {
	value, err := parseUint(str, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %v", str)
	}
	return uint8(value), nil
}

func parseValue16(str string) (uint16, error) {
	value, err := parseUint(str, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %v", str)
	}
	return uint16(value), nil
}

func formatValue(v uint8) string {
	return fmt.Sprintf("$%02x +%d", v, v)
}

func formatValue16(v uint16) string {
	return fmt.Sprintf("$%04x +%d", v, v)
}

func Dump(m memory.Memory, start uint16, end uint16, decode CharDecoder) string {
	var buf bytes.Buffer
	var chars bytes.Buffer

	a0 := start / 0x10 * 0x10
	a1 := end / 0x10 * 0x10
	if a1 != end {
		a1 += 0x10
	}
	for addr := a0; addr < a1; addr++ {
		if addr%0x10 == 0 {
			buf.WriteString(fmt.Sprintf("$%04x", addr))
			chars.Reset()
		}
		if addr < start || addr > end {
			buf.WriteString("   ")
			chars.WriteString(" ")
		} else {
			value := m.Load(addr)
			buf.WriteString(fmt.Sprintf(" %02x", value))
			ch, printable := decode(value)
			if printable {
				chars.WriteString(fmt.Sprintf("%c", ch))
			} else {
				chars.WriteString(".")
			}
		}
		if addr%0x10 == 7 {
			buf.WriteString(" ")
		}
		if addr%0x10 == 0x0f {
			buf.WriteString(" " + chars.String())
			if addr < end-1 {
				buf.WriteString("\n")
			}
		}
	}
	return buf.String()
}

var AsciiDecoder = func(code uint8) (rune, bool) {
	printable := code >= 32 && code < 128
	return rune(code), printable
}

var helpList = `
b   breakpoints
d   disassemble code
f   fill memory
g   go
h   halt
m   memory view
n   next
p   poke/peek memory
r   registers
s   step
si  state in
so  state out
t   trace
q   quit
`

var helpCmds = map[string]string{
	"b": `
Breakpoints

    b <address> {on|off}

Sets a breakpoint at <address> when using "on" and clears a breakpoint at
<address> when using "off". The CPU stops before executing address.
`,

	"d": `
Disassemble

    d [start-address [end-address]]

Disassemble code from [start-address] to [end-address] inclusive. If
[end-address] is not specified, disassemble an amount that can fit on a screen.
If [start-address] is not specified, use the current program counter as the
[start-address].
`,

	"f": `
Fill memory

    f <start-address> <end-address> <value>

Fill memory with <value> from <start-address> to <end-address> inclusive.
`,

	"g": `
Go

    g [address]

Go to [address] and start execution of the CPU there. If [address] is not
specified, use the current value of the program counter.
`,

	"h": `
Halt

    h

Halt execution of the CPU.
`,

	"m": `
Memory view

    m [start-address [end-address]]

Dump memory contents to the screen from [start-address] to [end-address]
inclusive. If [end-address] is not specified, show a full memory page.
If [start-address] is not specified, continue the dump from the last command.
`,

	"n": `
Next

    n

Disassemble the next instruction to execute.
`,

	"p": `
Peek/Poke

    p <address>

Peek at the memory contents at <address>. The value is displayed in the form
of $00 +000 with the hexadecimal value listed first followed by the decimal
value.

    p <address> <value>

Poke the memory at <address> with <value>.
`,

	"r": `
Registers

    r <name>

Display the value for the register with <name>.

    r <name> <value>

Set the <value> for register with <name>.
`,

	"s": `
Step

    s

Step through by executing the next instruction and then halting the CPU.
`,

	"so": `
State out

    so

Save the current machine state out to disk.
`,

	"si": `
State in

    si

Load the current machine state in from disk.
`,

	"t": `
Trace

    t

Toggle tracing of instructions executed by the CPU.
`,

	"q": `

Quit

    q[uit]

Quit to the operating system.
`,
}
