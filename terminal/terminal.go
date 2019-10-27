package terminal

import (
	"bufio"
	"errors"
	"os"
)

// Special keys
const (
	KeyNoSpl      = iota
	KeyArrowUp    // 27	91	65
	KeyArrowDown  // 27	91	66
	KeyArrowLeft  // 27	91	68
	KeyArrowRight // 27	91	67
	KeyPageUp     // 27 	91   5	126
	KeyPageDown   // 27	91	 6	126
	KeyHome       // 27	91	72
	KeyEnd        // 29	91	70
	KeyDelete     // 27	91	 3	126
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyIns
)

// matchSplKeys will try to match an input rune sequence
// to special keys and returns the number of speckal key
// combinations that are matched to the input sequence
func matchSplKeys(input []rune) (bool, int) {

outerLoop:
	for k, v := range specialKeys {
		if len(v) != len(input) {
			continue outerLoop
		}

		for j, r := range input {
			if r != v[j] {
				continue outerLoop
			}
		}
		return true, k
	}
	return false, 0
}

var specialKeys = map[int][]rune{
	KeyArrowUp:    []rune{27, 91, 65},
	KeyArrowDown:  []rune{27, 91, 66},
	KeyArrowLeft:  []rune{27, 91, 68},
	KeyArrowRight: []rune{27, 91, 67},
	KeyPageUp:     []rune{27, 91, 53, 126},
	KeyPageDown:   []rune{27, 91, 54, 126},
	KeyHome:       []rune{27, 91, 72},
	KeyEnd:        []rune{27, 91, 70},
	KeyDelete:     []rune{27, 91, 51, 126},
	KeyF1:         []rune{27, 79, 80},
	KeyF2:         []rune{27, 79, 81},
	KeyF3:         []rune{27, 79, 82},
	KeyF4:         []rune{27, 79, 83},
	KeyF5:         []rune{27, 91, 49, 53, 126},
	KeyF6:         []rune{27, 91, 49, 55, 126},
	KeyF7:         []rune{27, 91, 49, 56, 126},
	KeyF8:         []rune{27, 91, 49, 57, 126},
	KeyF9:         []rune{27, 91, 50, 48, 126},
	KeyIns:        []rune{27, 91, 50, 126},
}

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

// ReadKey reads a key from Stdin processing it for VT100 sequences.
// Stdin should be put in raw mode with
// VT100 processing enabled prior to using RawReadKey. If terminal read is set to
// timeout mode and no key is pressed, then ErrNoInput will be returned
func ReadKey() (Key, error) {

	r, n, err := bufr.ReadRune()
	if err != nil {
		return Key{}, err
	}

	// this code handles situation where a timeout has been set
	// but no key was pressed
	if n == 0 && err == nil {
		return Key{}, ErrNoInput
	}

	if r != 27 {
		return Key{r, KeyNoSpl}, nil
	}

	// nothing has been buffered, probably plain escape
	if bufr.Buffered() == 0 {
		return Key{27, KeyNoSpl}, nil
	}

	stack := []rune{27}
	for j := 0; j < 6; j++ {
		r, _, err := bufr.ReadRune()
		if err != nil {
			return Key{}, err
		}
		stack = append(stack, r)

		if match, key := matchSplKeys(stack); match {
			return Key{0, key}, nil
		}
	}
	return Key{0, KeyNoSpl}, nil
	// we have an escape key read. Let's check if it's a key we understand
}

// RawReadKey reads a key from Stdin. Stdin should be put in raw mode with
// VT100 processing enabled prior to using RawReadKey. If terminal read is set to
// timeout mode and no key is pressed, then ErrNoInput will be returned
func RawReadKey() (rune, error) {

	r, n, err := bufr.ReadRune()

	if err != nil {
		return 0, err
	}

	// this code handles situation where a timeout has been set
	// but no key was pressed
	if n == 0 && err == nil {
		return 0, ErrNoInput
	}

	return r, nil

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
