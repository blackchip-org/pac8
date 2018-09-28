package expect

import (
	"fmt"
	"reflect"
)

type Tester interface {
	Fatal(args ...interface{})
	Helper()
}

type WithClause struct {
	t      Tester
	format string
}

func With(t Tester) *WithClause {
	return &WithClause{t: t, format: "%v"}
}

func WithFormat(t Tester, format string) *WithClause {
	return &WithClause{t, format}
}

type ExpectClause struct {
	t      Tester
	format string
	have   interface{}
}

func (w WithClause) Expect(have interface{}) *ExpectClause {
	return &ExpectClause{t: w.t, format: w.format, have: have}
}

func (e ExpectClause) ToBe(want interface{}) {
	e.t.Helper()
	have := e.have
	haveType := reflect.TypeOf(have)
	wantType := reflect.TypeOf(want)
	if haveType != wantType {
		format := fmt.Sprintf("\n have: %v(%v) \n want: %v(%v)",
			haveType, e.format, wantType, e.format)
		message := fmt.Sprintf(format, have, want)
		e.t.Fatal(message)
		return
	}
	if !reflect.DeepEqual(have, want) {
		format := fmt.Sprintf("\n have: %v \n want: %v", e.format, e.format)
		message := fmt.Sprintf(format, have, want)
		e.t.Fatal(message)
		return
	}
}

func (e ExpectClause) ToPanic() {
	e.t.Helper()
	fn, ok := e.have.(func())
	if !ok {
		e.t.Fatal("value in expect should be a function")
	}
	defer func() {
		e.t.Helper()
		if r := recover(); r == nil {
			e.t.Fatal("expected panic")
		}
	}()
	fn()
}
