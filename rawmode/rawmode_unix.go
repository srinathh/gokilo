// +build !windows

package rawmode

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"golang.org/x/sys/unix"
)

// GetWindowSize returns the number of rows and columns in that order
// in the current console
func GetWindowSize() (int, int, error) {

	ws, err := unix.IoctlGetWinsize(unix.Stdout, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, fmt.Errorf("error fetching window size: %w", err)
	}
	if ws.Row == 0 || ws.Col == 0 {
		return 0, 0, fmt.Errorf("Got a zero size column or row")
	}

	return int(ws.Row), int(ws.Col), nil

}

// Enable switches the console from cooked or canonical mode to raw mode.
// It returns the current terminal settings for use in restoring console
// serlialized to a platform independent byte slice via gob
func Enable() ([]byte, error) {

	// Gets TermIOS data structure. From glibc, we find the cmd should be TCGETS
	// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/tcgetattr.c.html
	termios, err := unix.IoctlGetTermios(unix.Stdin, unix.TCGETS)
	if err != nil {
		return nil, fmt.Errorf("error fetching existing console settings: %w", err)
	}

	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(termios); err != nil {
		return nil, fmt.Errorf("error serializing existing console settings: %w", err)
	}

	// turn off echo & canonical mode by using a bitwise clear operator &^
	termios.Lflag = termios.Lflag &^ (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termios.Iflag = termios.Iflag &^ (unix.IXON | unix.ICRNL | unix.BRKINT | unix.INPCK | unix.ISTRIP)
	termios.Oflag = termios.Oflag &^ (unix.OPOST)
	termios.Cflag = termios.Cflag | unix.CS8
	//termios.Cc[unix.VMIN] = 0
	//termios.Cc[unix.VTIME] = 1
	// from the code of tcsetattr in glibc, we find that for TCSAFLUSH,
	// the corresponding command is TCSETSF
	// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/tcsetattr.c.html
	if err := unix.IoctlSetTermios(unix.Stdin, unix.TCSETSF, termios); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

// Restore restoes the console to a previous row setting
func Restore(original []byte) error {

	var termios unix.Termios

	if err := gob.NewDecoder(bytes.NewReader(original)).Decode(&termios); err != nil {
		return fmt.Errorf("error decoding terminal settings: %w", err)
	}

	if err := unix.IoctlSetTermios(unix.Stdin, unix.TCSETSF, &termios); err != nil {
		return fmt.Errorf("error restoring original console settings: %w", err)
	}
	return nil
}
