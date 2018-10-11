package mach

import (
	"testing"

	. "github.com/blackchip-org/pac8/util/expect"
)

func TestFitInWindow(t *testing.T) {
	frame := FitInWindow(1024, 768, 224, 288)
	With(t).Expect(frame.Scale).ToBe(int32(2))
}
