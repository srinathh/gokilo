package main

import "os"

// golang syscall main package is deprecated and
// points to sys/<os> packages to be used instead

const kiloVersion = "0.0.1"

func ctrlKey(b byte) int {
	return int(b & 0x1f)
}

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
