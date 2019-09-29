package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func safeExit(err error) {
	fmt.Fprint(os.Stdout, "\x1b[2J")
	fmt.Fprint(os.Stdout, "\x1b[H")

	if err1 := disableRawMode(); err1 != nil {
		fmt.Fprintf(os.Stderr, "Error: diabling raw mode: %s\r\n", err)
	}

	if err == nil {
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Error: %s\r\n", err)
	os.Exit(1)
}

// single space buffer to reduce allocations
var keyBuf = []byte{0}
var seq = []byte{0, 0, 0}
var errNoInput = errors.New("no input")

func rawReadKey() (byte, error) {
	n, err := os.Stdin.Read(keyBuf)
	switch {
	case err == io.EOF:
		return 0, errNoInput
	case err != nil:
		return 0, err
	case n == 0:
		return 0, errNoInput
	default:
		return keyBuf[0], nil
	}
}

func editorReadKey() int {

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
