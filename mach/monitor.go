package mach

import (
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

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
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

type Monitor struct {
	dasm    *cpu.Disassembler
	mach    *Mach
	cpu     cpu.CPU
	mem     memory.Memory
	in      io.ReadCloser
	out     *log.Logger
	rl      *readline.Instance
	lastCmd string
	memPtr  uint16
	dasmPtr uint16
}

func NewMonitor(mach *Mach) *Monitor {
	m := &Monitor{
		mach: mach,
		cpu:  mach.CPU,
		mem:  mach.Mem,
		in:   ioutil.NopCloser(os.Stdin),
		out:  log.New(os.Stdout, "", 0),
		dasm: mach.CPU.Info().NewDisassembler(mach.Mem),
	}
	mach.Callback = m.statusChange
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
			os.Exit(0)
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
		m.mach.Quit()
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
	go m.mach.Start()
	return nil
}

func (m *Monitor) halt(args []string) error {
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	m.mach.Stop()
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
		v := m.mem.Load(address)
		m.out.Printf("$%02x +%d\n", v, v)
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
	if err := checkLen(args, 0, 0); err != nil {
		return err
	}
	reason := ""
	if m.mach.Err != nil {
		reason = fmt.Sprintf(": %v", m.mach.Err)
	}
	m.out.Printf("[%v%v]\n", m.mach.Status, reason)
	m.out.Println(m.cpu.String())
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
	if m.mach.Tracing {
		m.mach.Trace(false)
	} else {
		m.mach.Trace(true)
	}
	return nil
}

func (m *Monitor) statusChange(s Status) {
	if s == Breakpoint {
		fmt.Println()
		m.registers([]string{})
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
