package terminal

import (
	"bufio"
	"errors"
	"os"
)

// Special keys
const (
	KeyNoSpl = iota
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyDelete
	KeyBackSpace
)

// single space buffer to reduce allocations
var keyBuf = []byte{0}
var seq = []byte{0, 0, 0}

// ErrNoInput indicates that there is no input when reading from keyboard
// in raw mode. This happens when timeout is set to a low number
var ErrNoInput = errors.New("no input")

// Key represents the key entered by the user
type Key struct {
	Regular rune
	Special int
}

var bufr = bufio.NewReader(os.Stdin)

// ReadKey reads a key from Stdin. Stdin should be put in raw mode with
// VT100 processing enabled prior to using RawReadKey. If terminal read is set to
// timeout mode and no key is pressed, then ErrNoInput will be returned
func ReadKey() (Key, error) {

	var ret Key

	r, n, err := bufr.ReadRune()

	if err != nil {
		return ret, err
	}

	// this code handles situation where a timeout has been set
	// but no key was pressed
	if n == 0 && err == nil {
		return ret, ErrNoInput
	}

	ret.Regular = r
	ret.Special = KeyNoSpl
	return ret, nil

}

/*

// ReadKey reads a key from Stdin using RawReadKey and interprets the VT100 sequences
// If the key pressed is a regular Unicode key,
// as either a regular
func ReadKey() int {

	for {
		key, err := rawReadKey()
		switch {
		case err == errNoInput:
			continue
		case err == io.EOF:
			safeExit(nil)
		case err != nil:
			safeExit(fmt.Errorf("Error reading key from STDIN: %s", err))
		case key == '\x1b':
			esc0, err := rawReadKey()
			if err == errNoInput || esc0 == '\x1b' {
				return '\x1b'
			}
			if err != nil {
				return 0
			}
			esc1, err := rawReadKey()
			if err == errNoInput {
				return '\x1b'
			}
			if err != nil {
				return 0
			}

			if esc0 == '[' {
				if esc1 >= '0' && esc1 <= '9' {
					esc2, err := rawReadKey()
					if err == errNoInput {
						return '\x1b'
					}
					if esc2 == '~' {
						switch esc1 {
						case '5':
							return keyPageUp
						case '6':
							return keyPageDown
						case '1', '7':
							return keyHome
						case '4', '8':
							return keyEnd
						case '3':
							return keyDelete
						}
					}

				} else {
					switch esc1 {
					case 'A':
						return keyArrowUp
					case 'B':
						return keyArrowDown
					case 'C':
						return keyArrowRight
					case 'D':
						return keyArrowLeft
					case 'H':
						return keyHome
					case 'F':
						return keyEnd
					}
				}
			} else if esc0 == 'O' {
				switch esc1 {
				case 'H':
					return keyHome
				case 'F':
					return keyEnd
				}
			}

		default:
			return int(key)
		}
	}
}

*/
