// +build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type winConsoleSettings struct {
	In  uint32
	Out uint32
}

func getWindowSize() (int, int, error) {

	info := windows.ConsoleScreenBufferInfo{}

	if err := windows.GetConsoleScreenBufferInfo(windows.Stdout, &info); err != nil {
		return 0, 0, fmt.Errorf("error fetching Screen SIze: %s", err)
	}

	return int(info.Window.Bottom - info.Window.Top + 1), int(info.Window.Right - info.Window.Left + 1), nil

}

func enableRawMode() error {

	oldSettings := winConsoleSettings{}

	if err := windows.GetConsoleMode(windows.Stdin, &oldSettings.In); err != nil {
		return fmt.Errorf("error gettingRaw mode: %s", err)
	}
	if err := windows.GetConsoleMode(windows.Stdout, &oldSettings.Out); err != nil {
		return fmt.Errorf("error gettingRaw mode: %s", err)
	}

	cfg.origTermCfg = oldSettings

	var inSettings uint32 = windows.ENABLE_EXTENDED_FLAGS | windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	var outSettings uint32 = windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_PROCESSED_OUTPUT | windows.DISABLE_NEWLINE_AUTO_RETURN

	if err := windows.SetConsoleMode(windows.Stdin, inSettings); err != nil {
		return fmt.Errorf("error setting Raw mode: %s", err)
	}

	if err := windows.SetConsoleMode(windows.Stdout, outSettings); err != nil {
		return fmt.Errorf("error setting output VT100: %s", err)
	}

	return nil
}

func disableRawMode() error {
	oldSettings := cfg.origTermCfg.(winConsoleSettings)
	if err := windows.SetConsoleMode(windows.Stdin, oldSettings.In); err != nil {
		return fmt.Errorf("error setting Raw mode: %s", err)
	}

	if err := windows.SetConsoleMode(windows.Stdout, oldSettings.Out); err != nil {
		return fmt.Errorf("error setting output VT100: %s", err)
	}
	return nil
}
