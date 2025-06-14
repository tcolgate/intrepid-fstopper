package main

// num is a number represented os hundredths
type num uint16

const (
	numMax num = (1 << 16) - 1
)

type numBuf [4]byte

// numOut renders n / 100 into out, rendered
// at most 1 decimal place
// - if >= 100(00) = " 123"
// - if <  100(00) = "99.9"
// - if <  100(00) = " 9.9"
func numOut(out *numBuf, n num) {
	hn := n / 1_00
	//ln := n % 1_00

	if hn >= 100 {
		// If we're over 100, just output as is
		out[0] = byte(' ')

		c1 := hn % 10
		c2 := (hn % 100) / 10
		c3 := hn / 100

		out[1] = byte('0' + c3)
		out[2] = byte('0' + c2)
		out[3] = byte('0' + c1)
		return
	}

	out[0] = byte('a')
	out[1] = byte('b')
	out[2] = byte('c')
	out[3] = byte('d')
	// output with decimal point
}
