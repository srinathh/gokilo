package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	syscall "golang.org/x/sys/unix"
)

// enableRawMode switches from cooked or canonical mode to raw mode
// by using syscalls. Currently this is the implrementation for Unix only
func enableRawMode() error {

	// Gets TermIOS data structure. From glibc, we find the cmd should be TCGETS
	// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/tcgetattr.c.html
	termios, err := syscall.IoctlGetTermios(syscall.Stdin, syscall.TCGETS)
	if err != nil {
		return err
	}

	cfg.origTermios = *termios

	// turn off echo & canonical mode by using a bitwise clear operator &^
	termios.Lflag = termios.Lflag &^ (syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	termios.Iflag = termios.Iflag &^ (syscall.IXON | syscall.ICRNL | syscall.BRKINT | syscall.INPCK | syscall.ISTRIP)
	termios.Oflag = termios.Oflag &^ (syscall.OPOST)
	termios.Cflag = termios.Cflag | syscall.CS8
	termios.Cc[syscall.VMIN] = 0
	termios.Cc[syscall.VTIME] = 1
	// from the code of tcsetattr in glibc, we find that for TCSAFLUSH,
	// the corresponding command is TCSETSF
	// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/tcsetattr.c.html
	if err := syscall.IoctlSetTermios(syscall.Stdin, syscall.TCSETSF, termios); err != nil {
		return err
	}

	return nil
}

func disableRawMode() error {
	if err := syscall.IoctlSetTermios(syscall.Stdin, syscall.TCSETSF, &cfg.origTermios); err != nil {
		return err
	}
	return nil
}

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
			if err == errNoInput {
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

func getWindowSize() (int, int, error) {

	ws, err := syscall.IoctlGetWinsize(syscall.Stdout, syscall.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	if ws.Row == 0 || ws.Col == 0 {
		return 0, 0, errors.New("got non zero column or row")
	}

	return int(ws.Row), int(ws.Col), nil

}
