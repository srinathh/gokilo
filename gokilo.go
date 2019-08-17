package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	// golang syscall main package is deprecated and
	// points to sys/<os> packages to be used instead
	syscall "golang.org/x/sys/unix"
)

/*** defines ***/

const kiloVersion = "0.0.1"

func ctrlKey(b byte) int {
	return int(b & 0x1f)
}

const (
	keyArrowUp    = 1000
	keyArrowDown  = 1001
	keyArrowLeft  = 1002
	keyArrowRight = 1003
	keyPageUp     = 1004
	keyPageDown   = 1005
)

/*** data ***/

type editorConfig struct {
	cx, cy      int
	screenRows  int
	screenCols  int
	origTermios syscall.Termios
}

var cfg editorConfig

/*** terminal ***/

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

func editorReadKey() (int, error) {

	for {
		key, err := rawReadKey()
		switch {
		case err == errNoInput:
			continue
		case err != nil:
			return 0, err
		case key == '\x1b':
			esc0, err := rawReadKey()
			if err == errNoInput {
				return '\x1b', nil
			}
			if err != nil {
				return 0, err
			}
			esc1, err := rawReadKey()
			if err == errNoInput {
				return '\x1b', nil
			}
			if err != nil {
				return 0, err
			}

			if esc0 == '[' {
				if esc1 >= '0' && esc1 <= '9' {
					esc2, err := rawReadKey()
					if err == errNoInput {
						return '\x1b', nil
					}
					if esc2 == '~' {
						switch {
						case esc1 == '5':
							return keyPageUp, nil
						case esc1 == '6':
							return keyPageDown, nil
						}
					}

				} else {
					switch {
					case esc1 == 'A':
						return keyArrowUp, nil
					case esc1 == 'B':
						return keyArrowDown, nil
					case esc1 == 'C':
						return keyArrowRight, nil
					case esc1 == 'D':
						return keyArrowLeft, nil
					}
				}
			}

		default:
			return int(key), nil
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

/*** output ***/
func editorRefreshScreen() {
	// clear screen
	ab := bytes.Buffer{}

	// hide cursor
	fmt.Fprint(&ab, "\x1b[?25l")

	// clear screen
	// fmt.Fprint(&ab, "\x1b[2J")

	// move cursor to top left
	fmt.Fprint(&ab, "\x1b[H")

	editorDrawRows(&ab)

	// reposition cursor
	//fmt.Fprint(&ab, "\x1b[H")
	fmt.Fprintf(&ab, "\x1b[%d;%dH", cfg.cy+1, cfg.cx+1)

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

func editorDrawRows(ab *bytes.Buffer) {
	for j := 0; j < cfg.screenRows; j++ {

		if j == cfg.screenRows/3 {
			welcomeMsg := fmt.Sprintf("Kilo Editor -- version %s", kiloVersion)
			welcomeLen := len(welcomeMsg)

			// if the message is too long to fit, truncate
			if welcomeLen > cfg.screenCols {
				welcomeMsg = welcomeMsg[:cfg.screenCols]
				welcomeLen = cfg.screenCols
			}
			padding := (cfg.screenCols - welcomeLen) / 2

			// if there is at least 1 padding required, use the Tilde to start line
			if padding > 0 {
				fmt.Fprint(ab, "~")
				padding--
			}

			// add appropriate number of spaces
			for i := 0; i < padding; i++ {
				fmt.Fprint(ab, " ")
			}
			fmt.Fprint(ab, welcomeMsg)

		} else {
			fmt.Fprint(ab, "~")
		}

		// clear to end of line
		fmt.Fprint(ab, "\x1b[K")

		if j < cfg.screenRows-1 {
			fmt.Fprint(ab, "\r\n")
		}
	}
}

/*** Input ***/

func editorProcessKeypress() error {

	b, err := editorReadKey()
	if err != nil {
		return err
	}

	switch b {
	case ctrlKey('q'):
		safeExit(nil)
	case keyArrowDown, keyArrowLeft, keyArrowRight, keyArrowUp:
		editorMoveCursor(b)
	case keyPageUp:
		for j := 0; j < cfg.screenRows; j++ {
			editorMoveCursor(keyArrowUp)
		}
	case keyPageDown:
		for j := 0; j < cfg.screenRows; j++ {
			editorMoveCursor(keyArrowDown)
		}
	}
	return nil
}

func editorMoveCursor(key int) {
	switch key {
	case keyArrowLeft:
		if cfg.cx != 0 {
			cfg.cx--
		}
	case keyArrowRight:
		if cfg.cx != cfg.screenCols-1 {
			cfg.cx++
		}
	case keyArrowDown:
		if cfg.cy != cfg.screenRows-1 {
			cfg.cy++
		}
	case keyArrowUp:
		if cfg.cy != 0 {
			cfg.cy--
		}
	}
}

/*** init ***/

func initEditor() error {
	rows, cols, err := getWindowSize()
	if err != nil {
		return err
	}
	cfg.screenRows = rows
	cfg.screenCols = cols
	return nil
}

func main() {

	if err := enableRawMode(); err != nil {
		safeExit(err)
	}

	if err := initEditor(); err != nil {
		safeExit(err)
	}

	for {
		editorRefreshScreen()
		if err := editorProcessKeypress(); err != nil {
			safeExit(err)
		}
	}
}
