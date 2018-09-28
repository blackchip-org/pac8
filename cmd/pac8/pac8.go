package main

import "github.com/blackchip-org/pac8/cabs"

func main() {
	mach := cabs.NewPacman().Mach()
	mach.Run()
}
