package z80

import "github.com/blackchip-org/pac8/bits"

func (cpu *CPU) storeIndImm(v uint8)    { cpu.mem.Store(cpu.fetch16(), v) }
func (cpu *CPU) store16IndImm(v uint16) { cpu.mem16.Store(cpu.fetch16(), v) }

func (cpu *CPU) storeA(v uint8) { cpu.A = v }
func (cpu *CPU) storeF(v uint8) { cpu.F = v }
func (cpu *CPU) storeB(v uint8) { cpu.B = v }
func (cpu *CPU) storeC(v uint8) { cpu.C = v }
func (cpu *CPU) storeD(v uint8) { cpu.D = v }
func (cpu *CPU) storeE(v uint8) { cpu.E = v }
func (cpu *CPU) storeH(v uint8) { cpu.H = v }
func (cpu *CPU) storeL(v uint8) { cpu.L = v }

func (cpu *CPU) storeAF(v uint16) { cpu.A, cpu.F = bits.Split(v) }
func (cpu *CPU) storeBC(v uint16) { cpu.B, cpu.C = bits.Split(v) }
func (cpu *CPU) storeDE(v uint16) { cpu.D, cpu.E = bits.Split(v) }
func (cpu *CPU) storeHL(v uint16) { cpu.H, cpu.L = bits.Split(v) }
func (cpu *CPU) storeSP(v uint16) { cpu.SP = v }

func (cpu *CPU) storeAF1(v uint16) { cpu.A1, cpu.F1 = bits.Split(v) }

func (cpu *CPU) storeIndBC(v uint8) { cpu.mem.Store(bits.Join(cpu.B, cpu.C), v) }
func (cpu *CPU) storeIndDE(v uint8) { cpu.mem.Store(bits.Join(cpu.D, cpu.E), v) }

func (cpu *CPU) loadImm() uint8       { return cpu.fetch() }
func (cpu *CPU) loadImm16() uint16    { return cpu.fetch16() }
func (cpu *CPU) loadIndImm() uint8    { return cpu.mem.Load(cpu.fetch16()) }
func (cpu *CPU) load16IndImm() uint16 { return cpu.mem16.Load(cpu.fetch16()) }

func (cpu *CPU) loadA() uint8 { return cpu.A }
func (cpu *CPU) loadF() uint8 { return cpu.F }
func (cpu *CPU) loadB() uint8 { return cpu.B }
func (cpu *CPU) loadC() uint8 { return cpu.C }
func (cpu *CPU) loadD() uint8 { return cpu.D }
func (cpu *CPU) loadE() uint8 { return cpu.E }
func (cpu *CPU) loadH() uint8 { return cpu.H }
func (cpu *CPU) loadL() uint8 { return cpu.L }

func (cpu *CPU) loadAF() uint16 { return bits.Join(cpu.A, cpu.F) }
func (cpu *CPU) loadBC() uint16 { return bits.Join(cpu.B, cpu.C) }
func (cpu *CPU) loadDE() uint16 { return bits.Join(cpu.D, cpu.E) }
func (cpu *CPU) loadHL() uint16 { return bits.Join(cpu.H, cpu.L) }
func (cpu *CPU) loadSP() uint16 { return cpu.SP }

func (cpu *CPU) loadAF1() uint16 { return bits.Join(cpu.A1, cpu.F1) }

func (cpu *CPU) loadIndBC() uint8 { return cpu.mem.Load(bits.Join(cpu.B, cpu.C)) }
func (cpu *CPU) loadIndDE() uint8 { return cpu.mem.Load(bits.Join(cpu.D, cpu.E)) }
