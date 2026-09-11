package telnet

import "bufio"

// Telnet protocol bytes (RFC 854).
const (
	SE   = 240 // end of subnegotiation
	NOP  = 241
	SB   = 250 // begin subnegotiation
	WILL = 251
	WONT = 252
	DO   = 253
	DONT = 254
	IAC  = 255 // interpret as command
)

// options
const (
	optEcho = 1
	optNAWS = 31
)

// iacFilter parser states.
type filterState int

const (
	stateData      filterState = iota // passing bytes through
	stateIAC                          // saw IAC, next byte is the command
	stateOption                       // saw WILL/WONT/DO/DONT, consume the option byte
	stateSubneg                       // inside SB ... SE, discard everything
	stateSubnegIAC                    // saw IAC inside a subnegotiation
)

type iacFilter struct {
	src   *bufio.Reader
	state filterState
}

func (f *iacFilter) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		if n > 0 && f.src.Buffered() == 0 {
			break // we have data, don't block waiting for more
		}
		b, err := f.src.ReadByte()
		if err != nil {
			return n, err
		}
		switch f.state {
		case stateData:
			if b == IAC {
				f.state = stateIAC
			} else if b != 0 { // drop NUL from CR NUL
				p[n] = b
				n++
			}
		case stateIAC:
			switch {
			case b == IAC: // escaped literal 0xFF
				p[n] = b
				n++
				f.state = stateData
			case b >= WILL && b <= DONT:
				f.state = stateOption // one option byte follows
			case b == SB:
				f.state = stateSubneg
			default: // NOP, GA, AYT, etc: single byte, nothing follows
				f.state = stateData
			}
		case stateOption:
			f.state = stateData
		case stateSubneg:
			if b == IAC {
				f.state = stateSubnegIAC
			}
		case stateSubnegIAC:
			if b == SE {
				f.state = stateData
			} else if b != IAC {
				f.state = stateSubneg
			}
		}
	}
	return n, nil
}
