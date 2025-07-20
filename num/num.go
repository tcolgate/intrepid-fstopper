// Package num implements operations on numbers encoded as 100ths.
// I wouldn't normally dedicate a package to this, but tinygo cannot
// build the test package without blowing out the flash size
package num

// Num is a number represented os hundredths
type Num uint32

const (
	Max    Num = (1 << 16) - 1
	IntMax Num = 9999
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

func IntOut(out *NumBuf, n Num) {
	cs := [4]Num{
		n / 1000,
		(n % 1000) / 100,
		(n % 100) / 10,
		n % 10,
	}

	ss := 4 - IntLen(n)
	for i := 0; i <= 3; i++ {
		if ss != 0 {
			out[i] = byte(' ')
			ss--
			continue
		}
		out[i] = byte('0' + cs[i])
	}
}

// IntOutLeft prints a number, as per the Out rules, but left justified
func IntOutLeft(out *NumBuf, n Num) {
	cs := [4]Num{
		n / 1000,
		(n % 1000) / 100,
		(n % 100) / 10,
		n % 10,
	}

	s := 4 - IntLen(n)
	for i := 0; i <= 3; i++ {
		if s >= 4 {
			out[i] = byte(' ')
			continue
		}

		out[i] = byte('0' + cs[s])
		s++
	}
}

// IntLen returns the printed length of the number, excluding any padding,
// this is currently 3 or 4
func IntLen(n Num) int {
	switch {
	case n >= 1000:
		return 4
	case n >= 100:
		return 3
	case n >= 10:
		return 2
	default:
		return 1
	}
}

func Mul(a Num, b int64) Num {
	switch b {
	case 0:
		return 0
	case 100:
		return a
	default:
		return Num((int64(a) * int64(b)) / 100)
	}
}

func Div(a Num, b int64) Num {
	switch b {
	case 0:
		// should panic, but then what?
		return Max
	case 100:
		return a
	default:
		return Num(((int64(a) * 100) / int64(b)))
	}
}
