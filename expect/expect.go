package expect

import (
	"fmt"
	"reflect"
)

// Tester defines the functions in the testing.T struct that are used
// by this package.
type Tester interface {
	Fatal(args ...interface{})
	Helper()
}

type WithClause struct {
	t      Tester
	format string
}

// With starts an assertion. Use *testing.T for Tester.
func With(t Tester) *WithClause {
	return &WithClause{t: t, format: "%v"}
}

// With starts an assertion with the given print format. The print format
// must contain a single verb which is used to print the both the want and
// have value. Use *testing.T for Tester.
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

// ToBe fails the test if the want value is not deeply equal to the have value
// used in the ExpectClause. If have and want are not the same type, have
// is converted to want if able, otherwise the test fails.
func (e ExpectClause) toBe(want interface{}, equal bool) {
	e.t.Helper()
	have := e.have
	haveType := reflect.TypeOf(have)
	wantType := reflect.TypeOf(want)

	if haveType.ConvertibleTo(wantType) {
		have = reflect.ValueOf(have).Convert(wantType).Interface()
	} else {
		if haveType != wantType {
			format := fmt.Sprintf("\n have: %v(%v) \n want: %v(%v)",
				haveType, e.format, wantType, e.format)
			message := fmt.Sprintf(format, have, want)
			e.t.Fatal(message)
			return
		}
	}
	var ok bool
	if equal {
		ok = reflect.DeepEqual(have, want)
	} else {
		ok = !reflect.DeepEqual(have, want)
	}
	if !ok {
		not := ""
		if !equal {
			not = "not "
		}
		format := fmt.Sprintf("\n have: %v \n want: %v%v", e.format, not, e.format)
		message := fmt.Sprintf(format, have, want)
		e.t.Fatal(message)
		return
	}
}

func (e ExpectClause) ToBe(want interface{}) {
	e.toBe(want, true)
}

func (e ExpectClause) NotToBe(want interface{}) {
	e.toBe(want, false)
}

// ToPanic fails the test if the have value in the ExpectClause does not panic
// when invoked. If the have value is not a function, the test fails.
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
