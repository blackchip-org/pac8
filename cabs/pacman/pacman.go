package pacman

// https://www.lomont.org/Software/Games/PacMan/PacmanEmulation.pdf

import (
	"fmt"
	"os"

	"github.com/blackchip-org/pac8/cpu/z80"
	"github.com/blackchip-org/pac8/mach"
	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/util/bits"
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
	spriteCoords    [8]spriteCoord
	dipSwitches     uint8
	watchdogReset   uint8
}

type voice struct {
	acc      [5]uint8
	waveform uint8
	freq     [5]uint8
	vol      uint8
}

type spriteCoord struct {
	x uint8
	y uint8
}

func New(renderer *sdl.Renderer) *mach.Mach {
	cab := &Pacman{}

	// Load ROMs
	e := []error{}
	rom0 := memory.LoadROM(&e, "pacman/pacman.6e", "e87e059c5be45753f7e9f33dff851f16d6751181")
	rom1 := memory.LoadROM(&e, "pacman/pacman.6f", "674d3a7f00d8be5e38b1fdc208ebef5a92d38329")
	rom2 := memory.LoadROM(&e, "pacman/pacman.6h", "8e47e8c2c4d6117d174cdac150392042d3e0a881")
	rom3 := memory.LoadROM(&e, "pacman/pacman.6j", "d4a70d56bb01d27d094d73db8667ffb00ca69cb9")
	vroms := VideoROM{
		Tiles: memory.LoadROM(&e, "pacman/pacman.5e", "06ef227747a440831c9a3a613b76693d52a2f0a9"),
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
	cab.mem = memory.NewPageMapped([]memory.Memory{
		rom0, // $0000 - $0fff
		rom1, // $1000 - $1fff
		rom2, // $2000 - $2fff
		rom3, // $3000 - $3fff
		ram,  // $4000 - $4fff
		io,   // $5000 - $50ff
	})
	// Mask out the bit 15 address line that is missing in Pacman
	cab.mem = memory.NewAddrMasked(cab.mem, 0x7fff)

	cab.cpu = z80.New(cab.mem)
	m := mach.New(cab.cpu)
	mapRegisters(&cab.regs, io)

	// Port 0 gets set with the partial interrupt pointer to be set
	// by the interrupting device
	cab.cpu.Ports.WO(0, &cab.intSelect)

	video, err := NewVideo(renderer, cab.mem, vroms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to initialize video: %v\n", err)
		os.Exit(1)
	}
	m.Display = video

	bits.Set(&cab.regs.in1, 4, true)         // Board test switch disabled
	bits.Set(&cab.regs.dipSwitches, 1, true) // 1 coin per game
	bits.Set(&cab.regs.dipSwitches, 7, true) // Normal ghost names

	// VBLANK interrupt
	video.Callback = func() {
		if m.Status == mach.Run && cab.regs.interruptEnable != 0 {
			cab.cpu.INT(cab.intSelect)
		}
	}

	return m
}

func mapRegisters(r *registers, io memory.IO) {
	for i := 0; i <= 0x3f; i++ {
		io.RO(i, &r.in0)
	}
	io.WO(0x00, &r.interruptEnable)
	io.WO(0x01, &r.soundEnable)
	io.WO(0x02, &r.unknown0)
	io.RW(0x03, &r.flipScreen)
	io.RW(0x04, &r.player1Lamp)
	io.RW(0x05, &r.player2Lamp)
	io.RW(0x06, &r.coinLockout)
	io.RW(0x07, &r.coinCounter)
	for i := 0x40; i <= 0x7f; i++ {
		io.RO(i, &r.in1)
	}
	for i, v := 0x40, 0; v < 3; i, v = i+6, v+1 {
		io.WO(i+0, &r.voices[v].acc[0])
		io.WO(i+1, &r.voices[v].acc[1])
		io.WO(i+2, &r.voices[v].acc[2])
		io.WO(i+3, &r.voices[v].acc[3])
		io.WO(i+4, &r.voices[v].acc[4])
		io.WO(i+5, &r.voices[v].waveform)
	}
	for i, v := 0x50, 0; v < 3; i, v = i+6, v+1 {
		io.WO(i+0, &r.voices[v].freq[0])
		io.WO(i+1, &r.voices[v].freq[1])
		io.WO(i+2, &r.voices[v].freq[2])
		io.WO(i+3, &r.voices[v].freq[3])
		io.WO(i+4, &r.voices[v].freq[4])
		io.WO(i+5, &r.voices[v].vol)
	}
	for i, s := 0x60, 0; s < 8; i, s = i+2, s+1 {
		io.WO(i+0, &r.spriteCoords[s].x)
		io.WO(i+1, &r.spriteCoords[s].y)
	}
	for i := 0x80; i <= 0xbf; i++ {
		io.RO(i, &r.dipSwitches)
	}
	for i := 0xc0; i <= 0xff; i++ {
		io.WO(i, &r.watchdogReset)
	}
}
