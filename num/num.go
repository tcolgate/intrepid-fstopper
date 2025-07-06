// Package num implements operations on numbers encoded as 100ths.
// I wouldn't normally dedicate a package to this, but tinygo cannot
// build the test package without blowing out the flash size
package num

// Num is a number represented os hundredths
type Num uint16

const (
	Max Num = (1 << 16) - 1
)

type NumBuf [4]byte

// numOut renders n / 100 into out, rendered
// at most 1 decimal place
// - if >= 100(00) = " 123"
// - if <  100(00) = "99.9"
// - if <  100(00) = " 9.9"
func Out(out *NumBuf, n Num) {
	hn := n / 1_00
	ln := n % 1_00

	if hn >= 100 {
		// If we're over 100, just output as is
		out[0] = byte(' ')

		c1 := hn / 100
		c2 := (hn % 100) / 10
		c3 := hn % 10

		out[1] = byte('0' + c1)
		out[2] = byte('0' + c2)
		out[3] = byte('0' + c3)
		return
	}

	c1 := hn / 10
	c2 := hn % 10

	out[0] = byte(' ')
	if c1 > 0 {
		out[0] = byte('0' + c1)
	}
	out[1] = byte('0' + c2)
	out[2] = byte('.')
	out[3] = byte('0' + (ln / 10))
	// output with decimal point
}

// OutLeft prints a number, as per the Out rules, but left justified
func OutLeft(out *NumBuf, n Num) {
	hn := n / 1_00
	ln := n % 1_00

	if hn >= 100 {
		// If we're over 100, just output as is
		out[3] = byte(' ')

		c1 := hn / 100
		c2 := (hn % 100) / 10
		c3 := hn % 10

		out[0] = byte('0' + c1)
		out[1] = byte('0' + c2)
		out[2] = byte('0' + c3)
		return
	}

	c1 := hn / 10
	c2 := hn % 10

	offset := 0
	out[3] = byte(' ')
	if c1 > 0 {
		out[0] = byte('0' + c1)
	} else {
		offset = -1
	}
	out[1+offset] = byte('0' + c2)
	out[2+offset] = byte('.')
	out[3+offset] = byte('0' + (ln / 10))
	// output with decimal point
}

// Len returns the printed length of the number, excluding any padding,
// this is currently 3 or 4
func Len(n Num) int {
	hn := n / 1_00

	if hn >= 100 {
		return 3
	}

	c1 := hn / 10

	if c1 > 0 {
		return 4
	} else {
		return 3
	}
}
