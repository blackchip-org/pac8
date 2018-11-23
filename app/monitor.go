package app

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	CmdMemory      = "m"
	CmdNext        = "n"
	CmdPokePeek    = "p"
	CmdRegisters   = "r"
	CmdStep        = "s"
	CmdTrace       = "t"
	CmdQuit        = "q"
	CmdQuitLong    = "quit"
)

const (
	memPageLen  = 0x100
	dasmPageLen = 0x3f
	maxArgs     = 0x100
)

type CharDecoder func(uint8) (rune, bool)

type Monitor struct {
	dasm    *proc.Disassembler
	mach    *machine.Mach
	cpu     proc.CPU
	mem     memory.Memory
	in      io.ReadCloser
	out     *log.Logger
	rl      *readline.Instance
	lastCmd string
	memPtr  uint16
	dasmPtr uint16
}

func NewMonitor(mach *machine.Mach) *Monitor {
	m := &Monitor{
		mach: mach,
		cpu:  mach.CPU,
		mem:  mach.Mem,
		in:   ioutil.NopCloser(os.Stdin),
		out:  log.New(os.Stdout, "", 0),
		dasm: mach.CPU.Info().NewDisassembler(mach.Mem),
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
		Prompt:      "monitor> ",
		HistoryFile: filepath.Join(usr.HomeDir, ".pac8-history"),
		Stdin:       m.in,
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
	case CmdMemory:
		err = m.memory(args, m.mach.CharDecoder)
	case CmdNext:
		err = m.next(args)
	case CmdPokePeek:
		err = m.pokePeek(args)
	case CmdRegisters:
		err = m.registers(args)
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
		if _, exists := m.mach.Breakpoints[address]; exists {
			m.out.Println("breakpoint on")
		} else {
			m.out.Println("breakpoint off")
		}
		return nil
	}
	switch args[1] {
	case "on":
		m.mach.Breakpoints[address] = struct{}{}
	case "off":
		delete(m.mach.Breakpoints, address)
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
		reason := ""
		if m.mach.Err != nil {
			reason = fmt.Sprintf(": %v", m.mach.Err)
		}
		m.out.Printf("[%v%v]\n", m.mach.Status, reason)
		m.out.Println(m.cpu.String())
		return nil
	}

	name := strings.ToUpper(args[0])
	reg, ok := m.mach.CPU.Info().Registers[name]
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
		if s == machine.Breakpoint {
			fmt.Println()
			m.registers([]string{})
		}
	case machine.TraceEvent:
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
