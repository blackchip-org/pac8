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
	"strconv"
	"strings"

	"github.com/blackchip-org/pac8/cpu"
	"github.com/blackchip-org/pac8/memory"
	"github.com/chzyer/readline"
)

const (
	CmdDisassemble = "d"
	CmdGo          = "g"
	CmdHalt        = "h"
	CmdRegisters   = "r"
	CmdTrace       = "t"
	CmdQuit        = "q"
	CmdQuitLong    = "quit"
)

const (
	memPageLen  = 0x100
	dasmPageLen = 0x3f
)

type Monitor struct {
	dasm    *cpu.Disassembler
	format  cpu.CodeFormatter
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

func NewMonitor(m *Mach) *Monitor {
	return &Monitor{
		mach:   m,
		cpu:    m.CPU,
		mem:    m.Mem,
		in:     ioutil.NopCloser(os.Stdin),
		out:    log.New(os.Stdout, "", 0),
		dasm:   m.NewDisassembler(),
		format: m.CPU.CodeFormatter(),
	}
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
		if err != nil {
			return err
		}
		m.parse(line)
	}
}

func (m *Monitor) parse(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}
	fields := strings.Split(line, " ")

	if len(fields) == 0 {
		return
	}

	cmd := fields[0]
	args := fields[1:]
	var err error
	switch cmd {
	case CmdDisassemble:
		err = m.disassemble(args)
	case CmdGo:
		err = m.go_(args)
	case CmdHalt:
		err = m.halt(args)
	case CmdRegisters:
		err = m.registers(args)
	case CmdTrace:
		err = m.trace(args)
	case CmdQuit, CmdQuitLong:
		os.Exit(0)
	default:
		err = fmt.Errorf("unknown command: %v", cmd)
	}

	if err != nil {
		m.out.Println(err)
	} else {
		m.lastCmd = cmd
	}
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
		addrStart = addr - 1
	}
	addrEnd := addrStart + uint16(dasmPageLen)
	if len(args) > 1 {
		addr, err := parseAddress(args[1])
		if err != nil {
			return err
		}
		addrEnd = addr - 1
	}
	m.dasm.SetPC(addrStart)
	for m.dasm.PC() <= addrEnd {
		statement := m.dasm.Next()
		m.out.Println(m.format(statement))
	}
	m.dasmPtr = m.dasm.PC()
	return nil
}

func (m *Monitor) go_(args []string) error {
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
