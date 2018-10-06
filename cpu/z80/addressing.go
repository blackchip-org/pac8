package z80

import (
	"github.com/blackchip-org/pac8/memory"
	"github.com/blackchip-org/pac8/util/bits"
)

func (cpu *CPU) storeNil(v uint8) {}

func (cpu *CPU) storeIndImm(v uint8)    { cpu.mem.Store(cpu.fetch16(), v) }
func (cpu *CPU) store16IndImm(v uint16) { memory.StoreLE(cpu.mem, cpu.fetch16(), v) }

func (cpu *CPU) storeA(v uint8)   { cpu.A = v }
func (cpu *CPU) storeB(v uint8)   { cpu.B = v }
func (cpu *CPU) storeC(v uint8)   { cpu.C = v }
func (cpu *CPU) storeD(v uint8)   { cpu.D = v }
func (cpu *CPU) storeE(v uint8)   { cpu.E = v }
func (cpu *CPU) storeH(v uint8)   { cpu.H = v }
func (cpu *CPU) storeL(v uint8)   { cpu.L = v }
func (cpu *CPU) storeI(v uint8)   { cpu.I = v }
func (cpu *CPU) storeR(v uint8)   { cpu.R = v }
func (cpu *CPU) storeIXH(v uint8) { cpu.IXH = v }
func (cpu *CPU) storeIXL(v uint8) { cpu.IXL = v }
func (cpu *CPU) storeIYH(v uint8) { cpu.IYH = v }
func (cpu *CPU) storeIYL(v uint8) { cpu.IYL = v }

func (cpu *CPU) storeAF(v uint16) { cpu.A, cpu.F = bits.Split(v) }
func (cpu *CPU) storeBC(v uint16) { cpu.B, cpu.C = bits.Split(v) }
func (cpu *CPU) storeDE(v uint16) { cpu.D, cpu.E = bits.Split(v) }
func (cpu *CPU) storeHL(v uint16) { cpu.H, cpu.L = bits.Split(v) }
func (cpu *CPU) storeSP(v uint16) { cpu.SP = v }
func (cpu *CPU) storeIX(v uint16) { cpu.IXH, cpu.IXL = bits.Split(v) }
func (cpu *CPU) storeIY(v uint16) { cpu.IYH, cpu.IYL = bits.Split(v) }

func (cpu *CPU) store16IndSP(v uint16) { memory.StoreLE(cpu.mem, cpu.SP, v) }

func (cpu *CPU) storeAF1(v uint16) { cpu.A1, cpu.F1 = bits.Split(v) }
func (cpu *CPU) storeBC1(v uint16) { cpu.B1, cpu.C1 = bits.Split(v) }
func (cpu *CPU) storeDE1(v uint16) { cpu.D1, cpu.E1 = bits.Split(v) }
func (cpu *CPU) storeHL1(v uint16) { cpu.H1, cpu.L1 = bits.Split(v) }

func (cpu *CPU) storeIndHL(v uint8) { cpu.mem.Store(bits.Join(cpu.H, cpu.L), v) }

func (cpu *CPU) storeIndBC(v uint8) { cpu.mem.Store(bits.Join(cpu.B, cpu.C), v) }
func (cpu *CPU) storeIndDE(v uint8) { cpu.mem.Store(bits.Join(cpu.D, cpu.E), v) }

func (cpu *CPU) loadZero() uint8      { return 0 }
func (cpu *CPU) loadImm() uint8       { return cpu.fetch() }
func (cpu *CPU) loadImm16() uint16    { return cpu.fetch16() }
func (cpu *CPU) loadIndImm() uint8    { return cpu.mem.Load(cpu.fetch16()) }
func (cpu *CPU) load16IndImm() uint16 { return memory.LoadLE(cpu.mem, cpu.fetch16()) }

func (cpu *CPU) loadA() uint8    { return cpu.A }
func (cpu *CPU) loadB() uint8    { return cpu.B }
func (cpu *CPU) loadC() uint8    { return cpu.C }
func (cpu *CPU) loadD() uint8    { return cpu.D }
func (cpu *CPU) loadE() uint8    { return cpu.E }
func (cpu *CPU) loadH() uint8    { return cpu.H }
func (cpu *CPU) loadL() uint8    { return cpu.L }
func (cpu *CPU) loadI() uint8    { return cpu.I }
func (cpu *CPU) loadR() uint8    { return cpu.R }
func (cpu *CPU) loadIXL() uint8  { return cpu.IXL }
func (cpu *CPU) loadIXH() uint8  { return cpu.IXH }
func (cpu *CPU) loadIYL() uint8  { return cpu.IYL }
func (cpu *CPU) loadIYH() uint8  { return cpu.IYH }
func (cpu *CPU) loadIndC() uint8 { return cpu.Ports.Load(uint16(cpu.C)) }

func (cpu *CPU) loadAF() uint16      { return bits.Join(cpu.A, cpu.F) }
func (cpu *CPU) loadBC() uint16      { return bits.Join(cpu.B, cpu.C) }
func (cpu *CPU) loadDE() uint16      { return bits.Join(cpu.D, cpu.E) }
func (cpu *CPU) loadHL() uint16      { return bits.Join(cpu.H, cpu.L) }
func (cpu *CPU) loadSP() uint16      { return cpu.SP }
func (cpu *CPU) loadIX() uint16      { return bits.Join(cpu.IXH, cpu.IXL) }
func (cpu *CPU) loadIY() uint16      { return bits.Join(cpu.IYH, cpu.IYL) }
func (cpu *CPU) load16IndSP() uint16 { return memory.LoadLE(cpu.mem, cpu.SP) }

func (cpu *CPU) loadAF1() uint16 { return bits.Join(cpu.A1, cpu.F1) }
func (cpu *CPU) loadBC1() uint16 { return bits.Join(cpu.B1, cpu.C1) }
func (cpu *CPU) loadDE1() uint16 { return bits.Join(cpu.D1, cpu.E1) }
func (cpu *CPU) loadHL1() uint16 { return bits.Join(cpu.H1, cpu.L1) }

func (cpu *CPU) loadIndHL() uint8 { return cpu.mem.Load(bits.Join(cpu.H, cpu.L)) }

func (cpu *CPU) loadIndBC() uint8 { return cpu.mem.Load(bits.Join(cpu.B, cpu.C)) }
func (cpu *CPU) loadIndDE() uint8 { return cpu.mem.Load(bits.Join(cpu.D, cpu.E)) }

func (cpu *CPU) loadIndIX() uint8 {
	ix := bits.Join(cpu.IXH, cpu.IXL)
	cpu.iaddr = bits.Displace(ix, cpu.delta)
	return cpu.mem.Load(cpu.iaddr)
}

func (cpu *CPU) loadIndIY() uint8 {
	iy := bits.Join(cpu.IYH, cpu.IYL)
	cpu.iaddr = bits.Displace(iy, cpu.delta)
	return cpu.mem.Load(cpu.iaddr)
}

func (cpu *CPU) storeIndIX(v uint8) {
	ix := bits.Join(cpu.IXH, cpu.IXL)
	addr := bits.Displace(ix, cpu.delta)
	cpu.mem.Store(addr, v)
}

func (cpu *CPU) storeIndIY(v uint8) {
	iy := bits.Join(cpu.IYH, cpu.IYL)
	addr := bits.Displace(iy, cpu.delta)
	cpu.mem.Store(addr, v)
}

func (cpu *CPU) storeLastInd(v uint8) {
	cpu.mem.Store(cpu.iaddr, v)
}

func (cpu *CPU) outIndImm(v uint8) {
	addr := uint16(cpu.fetch())
	cpu.Ports.Store(addr, v)
}

func (cpu *CPU) inIndImm() uint8 {
	addr := uint16(cpu.fetch())
	return cpu.Ports.Load(addr)
}

func (cpu *CPU) outIndC(v uint8) {
	cpu.Ports.Store(uint16(cpu.C), v)
}

func (cpu *CPU) inIndC() uint8 {
	return cpu.Ports.Load(uint16(cpu.C))
}
