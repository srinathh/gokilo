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

	return int(info.Size.Y), int(info.Size.X), nil

}
