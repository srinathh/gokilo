// +build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func getWindowSize() (int, int, error) {

	info := windows.ConsoleScreenBufferInfo{}

	if err := windows.GetConsoleScreenBufferInfo(windows.Stdout, &info); err != nil {
		return 0, 0, fmt.Errorf("error fetching Screen SIze: %s", err)
	}

	return int(info.Window.Bottom - info.Window.Top + 1), int(info.Window.Right - info.Window.Left + 1), nil

}

func enableRawMode() error {

	var inSettings uint32 = windows.ENABLE_EXTENDED_FLAGS | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	var outSettings uint32 = windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_PROCESSED_OUTPUT

	if err := windows.SetConsoleMode(windows.Stdin, inSettings); err != nil {
		return fmt.Errorf("error setting Raw mode: %s", err)
	}

	if err := windows.SetConsoleMode(windows.Stdout, outSettings); err != nil {
		return fmt.Errorf("error setting output VT100: %s", err)
	}

	return nil
}

func disableRawMode() error {
	return nil
}
