// +build windows

package rawmode

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"golang.org/x/sys/windows"
)

type winConsoleSettings struct {
	In  uint32
	Out uint32
}

// GetWindowSize returns the number of rows and columns in that order
// in the current console
func GetWindowSize() (int, int, error) {

	info := windows.ConsoleScreenBufferInfo{}

	if err := windows.GetConsoleScreenBufferInfo(windows.Stdout, &info); err != nil {
		return 0, 0, fmt.Errorf("error fetching Screen Size: %w", err)
	}

	return int(info.Window.Bottom - info.Window.Top + 1), int(info.Window.Right - info.Window.Left + 1), nil

}

// Enable switches the console from cooked or canonical mode to raw mode.
// It returns the current terminal settings for use in restoring console
// serlialized to a platform independent byte slice via gob
func Enable() ([]byte, error) {

	oldSettings := winConsoleSettings{}

	if err := windows.GetConsoleMode(windows.Stdin, &oldSettings.In); err != nil {
		return nil, fmt.Errorf("error gettingRaw mode: %s", err)
	}
	if err := windows.GetConsoleMode(windows.Stdout, &oldSettings.Out); err != nil {
		return nil, fmt.Errorf("error gettingRaw mode: %w", err)
	}

	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(oldSettings); err != nil {
		return nil, fmt.Errorf("error serializing existing console settings: %w", err)
	}

	var inSettings uint32 = windows.ENABLE_EXTENDED_FLAGS | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	var outSettings uint32 = windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_PROCESSED_OUTPUT | windows.DISABLE_NEWLINE_AUTO_RETURN

	if err := windows.SetConsoleMode(windows.Stdin, inSettings); err != nil {
		return buf.Bytes(), fmt.Errorf("error setting raw mode: %w", err)
	}

	if err := windows.SetConsoleMode(windows.Stdout, outSettings); err != nil {
		return buf.Bytes(), fmt.Errorf("error setting output VT100: %s", err)
	}

	return buf.Bytes(), nil
}

// Restore restoes the console to a previous row setting
func Restore(original []byte) error {

	var winCS winConsoleSettings
	if err := gob.NewDecoder(bytes.NewReader(original)).Decode(&winCS); err != nil {
		return fmt.Errorf("error decoding terminal settings: %w", err)
	}

	if err := windows.SetConsoleMode(windows.Stdin, winCS.In); err != nil {
		return fmt.Errorf("error setting Raw mode: %w", err)
	}

	if err := windows.SetConsoleMode(windows.Stdout, winCS.Out); err != nil {
		return fmt.Errorf("error setting output VT100: %w", err)
	}
	return nil
}
