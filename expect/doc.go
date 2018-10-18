/*
Package expect is a simple assertion library.

Example:

	import (
		testing
	 	. "github.com/blackchip-org/pac8/expect"
	)

	func TestTwo(t *testing.T) {
		With(t).Expect(1 + 4).ToBe(2)
	}

When an assertion fails, the output will look like the following:

	--- FAIL: TestTwo (0.00s)
	/home/me/go/src/foo/my_test.go:7:
		 have: 5
		 want: 2

Format values in the assertion output:

	func TestHex(t *testing.T) {
		WithFormat(t, "$%04X").Expect(0x1234).ToBe(0xabcd)
	}

The failure for this test looks like:

	--- FAIL: TestHex (0.00s)
	/home/me/go/src/foo/my_test.go:7:
		 have: $1234
		 want: $ABCD
*/
package expect
