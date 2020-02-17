package terminal

import (
	"bufio"
	"fmt"
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

// ErrNoInput indicates that there is no input when reading from keyboard
// in raw mode. This happens when timeout is set to a low number
var ErrNoInput = fmt.Errorf("no input")

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
	// we couldn't make out the special key, let's just return escape
	// this is probably wrong but unless we have a custom bufio.Reader,
	// we can't do better
	return Key{27, KeyNoSpl}, nil
}
