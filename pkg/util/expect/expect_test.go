package expect

import (
	"testing"
)

func TestToBeEqual(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(2).ToBe(2)
	if dt.args != nil {
		t.Fatalf("expected to be equal")
	}
}

func TestNotToBeEqual(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(2).NotToBe(2)
	if dt.args == nil {
		t.Fatalf("expected to be not equal")
	}
	if dt.args[0] != "\n have: 2 \n want: not 2" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestToBeNotEqual(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(4).ToBe(2)
	if dt.args == nil {
		t.Fatalf("expected to be not equal")
	}
	if dt.args[0] != "\n have: 4 \n want: 2" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestNotToBeNotEqual(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(4).NotToBe(2)
	if dt.args != nil {
		t.Fatalf("expected to be not equal")
	}
}

func TestNotEqualFormat(t *testing.T) {
	dt := &debugTester{}
	WithFormat(dt, "$%04X").Expect(0x1234).ToBe(0xabcd)
	if dt.args == nil {
		t.Fatalf("expected to be not equal")
	}
	if dt.args[0] != "\n have: $1234 \n want: $ABCD" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestConvertable(t *testing.T) {
	dt := &debugTester{}
	a := int(4)
	b := int32(4)
	With(dt).Expect(a).ToBe(b)
	if dt.args != nil {
		t.Fatalf("expected to be equal")
	}
}

func TestNotConvertable(t *testing.T) {
	dt := &debugTester{}
	a := "4"
	b := 4
	With(dt).Expect(a).ToBe(b)
	if dt.args == nil {
		t.Fatalf("expected to be not equal")
	}
	if dt.args[0] != "\n have: string(4) \n want: int(4)" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestDeepEquals(t *testing.T) {
	dt := &debugTester{}
	a := [][]int{
		[]int{1, 2, 3},
		[]int{4},
	}
	b := [][]int{
		[]int{1, 2, 3},
		[]int{4},
	}
	With(dt).Expect(a).ToBe(b)
	if dt.args != nil {
		t.Fatalf("expected to be equal")
	}
}

func TestNotDeepEquals(t *testing.T) {
	dt := &debugTester{}
	a := [][]int{
		[]int{1, 2, 3},
		[]int{4},
		[]int{5},
	}
	b := [][]int{
		[]int{1, 2, 3},
		[]int{4},
	}
	With(dt).Expect(a).ToBe(b)
	if dt.args == nil {
		t.Fatalf("expected to be not equal")
	}
	if dt.args[0] != "\n have: [[1 2 3] [4] [5]] \n want: [[1 2 3] [4]]" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestPanic(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(func() { panic("panic") }).ToPanic()
	if dt.args != nil {
		t.Fatalf("expected panic")
	}
}

func TestNotPanic(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(func() {}).ToPanic()
	if dt.args == nil {
		t.Fatalf("expected not to panic")
	}
	if dt.args[0] != "expected panic" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestNotFuncPanic(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(2).ToPanic()
	if dt.args == nil {
		t.Fatalf("expected an error")
	}
	if dt.args[0] != "value in expect should be a function" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestExpectNil(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(nil).ToBe(nil)
	if dt.args != nil {
		t.Fatalf("expected to be equal")
	}
}

func TestExpectNilFail(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(nil).ToBe(2)
	if dt.args == nil {
		t.Fatalf("expected to be not equal")
	}
	if dt.args[0] != "\n have: <nil> \n want: 2" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

func TestExpectToBeNilFail(t *testing.T) {
	dt := &debugTester{}
	With(dt).Expect(2).ToBe(nil)
	if dt.args == nil {
		t.Fatalf("expected to be not equal")
	}
	if dt.args[0] != "\n have: 2 \n want: <nil>" {
		t.Fatalf("unexpected: %v", dt.args[0])
	}
}

type debugTester struct {
	args []interface{}
}

func (t *debugTester) Fatal(args ...interface{}) {
	t.args = args
}

func (t *debugTester) Helper() {}
