package z80

func (cpu *CPU) storeA(value uint8) { cpu.A = value }
func (cpu *CPU) storeF(value uint8) { cpu.F = value }
func (cpu *CPU) storeB(value uint8) { cpu.B = value }
func (cpu *CPU) storeC(value uint8) { cpu.C = value }
func (cpu *CPU) storeD(value uint8) { cpu.D = value }
func (cpu *CPU) storeE(value uint8) { cpu.E = value }
func (cpu *CPU) storeH(value uint8) { cpu.H = value }
func (cpu *CPU) storeL(value uint8) { cpu.L = value }

func (cpu *CPU) loadA() uint8 { return cpu.A }
func (cpu *CPU) loadF() uint8 { return cpu.F }
func (cpu *CPU) loadB() uint8 { return cpu.B }
func (cpu *CPU) loadC() uint8 { return cpu.C }
func (cpu *CPU) loadD() uint8 { return cpu.D }
func (cpu *CPU) loadE() uint8 { return cpu.E }
func (cpu *CPU) loadH() uint8 { return cpu.H }
func (cpu *CPU) loadL() uint8 { return cpu.L }
