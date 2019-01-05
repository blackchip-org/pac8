package galaga

import (
	"fmt"
	"time"

	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/component/memory"
	"github.com/blackchip-org/pac8/component/proc"
	"github.com/blackchip-org/pac8/component/proc/z80"
	"github.com/blackchip-org/pac8/machine"
	"github.com/blackchip-org/pac8/pkg/pac8"
	"github.com/blackchip-org/pac8/pkg/util/state"
)

type Galaga struct {
	spec *machine.Spec
	regs Registers
}

type Config struct {
	Name string
}

type Registers struct {
	InterruptEnable0 uint8 // low bit
	InterruptEnable1 uint8 // low bit
	InterruptEnable2 uint8 // low bit
	DipSwitches      [8]uint8
}

var codeSegments = []string{"code1", "code2", "code3"}

func New(env pac8.Env, config Config, roms memory.Set) (machine.System, error) {
	sys := &Galaga{}

	ram := memory.NewRAM(0x2000)
	io := memory.NewIO(0x100)
	xram := memory.NewRAM(0x1000)
	xram2 := memory.NewRAM(0x1000)

	mem := make([]memory.Memory, 4, 4)
	cpu := make([]*z80.CPU, 4, 4)
	for i := 0; i < 3; i++ {
		m := memory.NewBlockMapper()
		m.Map(0x0000, roms[codeSegments[i]])
		m.Map(0x6800, io)
		m.Map(0x7000, xram)
		m.Map(0x8000, ram)
		m.Map(0xa000, xram2)
		mem[i] = memory.NewPageMapped(m.Blocks)
		spy := memory.NewSpy(mem[i])
		mem[i] = spy
		cpu[i] = z80.New(mem[i])

		nCore := i
		coreCPU := cpu[i]
		spy.Callback(func(e memory.Event) {
			fmt.Printf("core %v at %04x: %v\n", nCore+1, coreCPU.PC(), e)
		})
		//spy.WatchW(0x7000)
		//spy.WatchW(0x7100)
	}
	mem[3] = mem[0]
	mem[0].Store(0x9100, 0xff)
	mem[0].Store(0x9101, 0xff)

	mapRegisters(&sys.regs, io)

	video, err := NewVideo(env.Renderer, mem[0], roms)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize video: %v", err)
	}

	bits.Set(&sys.regs.DipSwitches[3], 0, true)
	bits.Set(&sys.regs.DipSwitches[5], 0, true)
	bits.Set(&sys.regs.DipSwitches[6], 0, true)

	hackCPU := &HackCPU{cpu: cpu[0], mem: mem[0]}
	sys.spec = &machine.Spec{
		Name:        config.Name,
		CharDecoder: GalagaDecoder,
		CPU:         []proc.CPU{cpu[0], cpu[1], cpu[2], hackCPU},
		Mem:         mem,
		Display:     video,
		TickCallback: func(m *machine.Mach) {
			if m.Status != machine.Run {
				return
			}
			if sys.regs.InterruptEnable0 != 0 {
				cpu[0].INT(0)
				cpu[0].NMI()
			}
			if sys.regs.InterruptEnable1 != 0 {
				cpu[1].INT(0)
			}
			if sys.regs.InterruptEnable2 != 0 {
				cpu[2].INT(0)
			}
		},
		TickRate: time.Duration(16670 * time.Microsecond),
	}
	return sys, nil
}

func (g Galaga) Spec() *machine.Spec {
	return g.spec
}

func (g *Galaga) Save(enc *state.Encoder) {}

func (g *Galaga) Restore(dec *state.Decoder) {}

func mapRegisters(r *Registers, io memory.IO) {
	pm := memory.NewPortMapper(io)
	for i := 0; i <= 7; i++ {
		pm.RW(i, &r.DipSwitches[i])
	}
	pm.RW(0x20, &r.InterruptEnable0)
	pm.RW(0x21, &r.InterruptEnable1)
	pm.RW(0x22, &r.InterruptEnable2)
}

type HackCPU struct {
	cpu   proc.CPU
	mem   memory.Memory
	count int
	stuff bool
}

func (h HackCPU) PC() uint16 {
	return 0
}

func (h HackCPU) SetPC(_ uint16) {
}

func (h *HackCPU) Next() {
	if !h.stuff {
		if h.cpu.PC() >= 0x37ec && h.cpu.PC() <= 0x37f2 {
			h.stuff = true
		}
	}
	h.count++
	if h.stuff {
		h.mem.Store(0x7100, 0x10)
		v := uint8(h.count)
		h.mem.Store(0x92a0, v)
	}
}

func (h HackCPU) Ready() bool {
	return true
}

func (h HackCPU) Info() proc.Info {
	return proc.Info{
		NewDisassembler: func(memory.Memory) *proc.Disassembler {
			return nil
		},
	}
}

func (h HackCPU) Save(*state.Encoder) {}

func (h HackCPU) Restore(*state.Decoder) {}

func (h HackCPU) String() string {
	return ""
}
