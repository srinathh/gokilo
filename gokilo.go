package main

import (
	"os"
)

const kiloVersion = "0.0.1"

var cfg Config
var editor Editor

func ctrlKey(b byte) int {
	return int(b & 0x1f)
}

func initEditor() error {
	rows, cols, err := getWindowSize()
	if err != nil {
		return err
	}
	cfg.ScreenRows = rows
	cfg.ScreenRows = cfg.ScreenRows - 2
	cfg.ScreenCols = cols
	editor.QuitTimes = kiloQuitTimes
	editor.LastMatch = -1
	editor.Direction = 1
	editorSetStatusMsg("HELP: Ctrl + S to save | Ctrl + Q to exit | Ctrl + F to find")
	return nil
}

func main() {

	if err := enableRawMode(); err != nil {
		safeExit(err)
	}

	if err := initEditor(); err != nil {
		safeExit(err)
	}

	if len(os.Args) >= 2 {
		if err := editorOpen(os.Args[1]); err != nil {
			safeExit(err)
		}
	}

	for {
		editorRefreshScreen()
		if err := editorProcessKeypress(); err != nil {
			safeExit(err)
		}
	}
}
