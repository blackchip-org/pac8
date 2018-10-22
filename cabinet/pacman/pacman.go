package pacman

// https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

import (
	"fmt"
	"os"
	"time"

	"github.com/blackchip-org/pac8/app"
	"github.com/blackchip-org/pac8/bits"
	"github.com/blackchip-org/pac8/cpu/z80"
	"github.com/blackchip-org/pac8/mach"
	"github.com/blackchip-org/pac8/memory"
	"github.com/veandco/go-sdl2/sdl"
)

type Pacman struct {
	mem       memory.Memory
	cpu       *z80.CPU
	regs      registers
	intSelect uint8 // value sent during interrupt to select vector (port 0)
	tiles     *sdl.Texture
}

type registers struct {
	in0             uint8 // joystick and coin slot
	interruptEnable uint8
	soundEnable     uint8
	unknown0        uint8
	flipScreen      uint8
	player1Lamp     uint8
	player2Lamp     uint8
	coinLockout     uint8
	coinCounter     uint8
	in1             uint8 // joystick and start buttons
	voices          [3]voice
	dipSwitches     uint8
	watchdogReset   uint8
}

type voice struct {
	acc      [5]uint8
	waveform uint8
	freq     [5]uint8
	vol      uint8
}

func New(renderer *sdl.Renderer) *mach.Mach {
	cab := &Pacman{}

	// Load ROMs
	e := []error{}
	rom0 := app.LoadROM(&e, "pacman/pacman.6e", "e87e059c5be45753f7e9f33dff851f16d6751181")
	rom1 := app.LoadROM(&e, "pacman/pacman.6f", "674d3a7f00d8be5e38b1fdc208ebef5a92d38329")
	rom2 := app.LoadROM(&e, "pacman/pacman.6h", "8e47e8c2c4d6117d174cdac150392042d3e0a881")
	rom3 := app.LoadROM(&e, "pacman/pacman.6j", "d4a70d56bb01d27d094d73db8667ffb00ca69cb9")
	vroms := VideoROM{
		Tiles:   app.LoadROM(&e, "pacman/pacman.5e", "06ef227747a440831c9a3a613b76693d52a2f0a9"),
		Sprites: app.LoadROM(&e, "pacman/pacman.5f", "4a937ac02216ea8c96477d4a15522070507fb599"),
		Color:   app.LoadROM(&e, "pacman/82s123.7f", "8d0268dee78e47c712202b0ec4f1f51109b1f2a5"),
		Palette: app.LoadROM(&e, "pacman/82s126.4a", "19097b5f60d1030f8b82d9f1d3a241f93e5c75d6"),
	}

	// Any errors while loading ROMs?
	for _, err := range e {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	if len(e) > 0 {
		os.Exit(1)
	}

	ram := memory.NewRAM(0x1000)
	io := memory.NewIO(0x100)
	cab.mem = memory.NewPageMapped([]memory.Block{
		memory.NewBlock(0x0000, rom0),
		memory.NewBlock(0x1000, rom1),
		memory.NewBlock(0x2000, rom2),
		memory.NewBlock(0x3000, rom3),
		memory.NewBlock(0x4000, ram),
		memory.NewBlock(0x5000, io),
	})
	// Mask out the bit 15 address line that is missing in Pacman
	cab.mem = memory.NewMasked(cab.mem, 0x7fff)

	cab.cpu = z80.New(cab.mem)
	m := mach.New(cab.mem, cab.cpu)

	video, err := NewVideo(renderer, cab.mem, vroms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize video: %v\n", err)
		os.Exit(1)
	}
	m.Display = video
	keyboard := NewKeyboard(&cab.regs)
	m.UI = keyboard

	mapRegisters(&cab.regs, io, video)

	// Port 0 gets set with the partial interrupt pointer to be set
	// by the interrupting device
	cab.cpu.Map.WO(0, &cab.intSelect)

	bits.Set(&cab.regs.in1, 4, true)          // Board test switch disabled
	bits.Set(&cab.regs.dipSwitches, 0, true)  // 1 coin per game
	bits.Set(&cab.regs.dipSwitches, 1, false) // ...
	bits.Set(&cab.regs.dipSwitches, 3, true)  // 3 lives
	bits.Set(&cab.regs.dipSwitches, 7, true)  // Normal ghost names

	// Different type of crash when these are set
	cab.regs.in0 = 0x3f
	cab.regs.in1 = 0x7f

	// 16.67 milliseconds for VBLANK interrupt
	m.TickRate = time.Duration(16670 * time.Microsecond)
	video.Callback = func() {
		if m.Status == mach.Run && cab.regs.interruptEnable != 0 {
			cab.cpu.INT(cab.intSelect)
		}
	}

	return m
}

func mapRegisters(r *registers, io memory.IO, v *Video) {
	pm := memory.NewPortMapper(io)
	for i := 0; i <= 0x3f; i++ {
		pm.RO(i, &r.in0)
	}
	pm.WO(0x00, &r.interruptEnable)
	pm.WO(0x01, &r.soundEnable)
	pm.WO(0x02, &r.unknown0)
	pm.RW(0x03, &r.flipScreen)
	pm.RW(0x04, &r.player1Lamp)
	pm.RW(0x05, &r.player2Lamp)
	pm.RW(0x06, &r.coinLockout)
	pm.RW(0x07, &r.coinCounter)
	for i := 0x40; i <= 0x7f; i++ {
		pm.RO(i, &r.in1)
	}
	for i, v := 0x40, 0; v < 3; i, v = i+6, v+1 {
		pm.WO(i+0, &r.voices[v].acc[0])
		pm.WO(i+1, &r.voices[v].acc[1])
		pm.WO(i+2, &r.voices[v].acc[2])
		pm.WO(i+3, &r.voices[v].acc[3])
		pm.WO(i+4, &r.voices[v].acc[4])
		pm.WO(i+5, &r.voices[v].waveform)
	}
	for i, v := 0x50, 0; v < 3; i, v = i+6, v+1 {
		pm.WO(i+0, &r.voices[v].freq[0])
		pm.WO(i+1, &r.voices[v].freq[1])
		pm.WO(i+2, &r.voices[v].freq[2])
		pm.WO(i+3, &r.voices[v].freq[3])
		pm.WO(i+4, &r.voices[v].freq[4])
		pm.WO(i+5, &r.voices[v].vol)
	}
	for i, s := 0x60, 0; s < 8; i, s = i+2, s+1 {
		pm.WO(i+0, &v.spriteCoords[s].x)
		pm.WO(i+1, &v.spriteCoords[s].y)
	}
	for i := 0x80; i <= 0xbf; i++ {
		pm.RO(i, &r.dipSwitches)
	}
	for i := 0xc0; i <= 0xff; i++ {
		pm.WO(i, &r.watchdogReset)
	}
}
