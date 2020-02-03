package main

import (
	"flag"
	"fmt"
	"gokilo/rawmode"
	"gokilo/terminal"
	"os"
)

func ctrlKey(b byte) rune {
	return rune(b & 0x1f)
}

const kiloVersion = "0.0.2"

func safeExit(origCfg []byte, err error) {
	fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")

	if err1 := rawmode.Restore(origCfg); err1 != nil {
		fmt.Fprintf(os.Stderr, "Error: disabling raw mode: %s\r\n", err)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\r\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// SafeExit is a global function that can be called to exit safely
var SafeExit func(error)

func main() {

	origCfg, err := rawmode.Enable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error enabling raw mode: %v", err)
		os.Exit(1)
	}

	SafeExit = func(error) { safeExit(origCfg, err) }

	rows, cols, err := rawmode.GetWindowSize()
	if err != nil {
		SafeExit(fmt.Errorf("couldn't get window size: %v", err))
	}

	view := NewView(rows, cols)

	var e *Editor
	//var s *Session

	flag.Parse()
	filename := flag.Arg(0)

	if flag.Arg(0) == "" {
		e = NewEditor()
	} else {
		e, err = NewEditorFromFile(filename)
		if err != nil {
			SafeExit(fmt.Errorf("couldn't open file %s: %v", filename, err))
		}
	}
	//s = NewSession(filename)

	for {
		view.RefreshScreen(e)

		k, err := terminal.ReadKey()
		if err != nil {
			SafeExit(fmt.Errorf("Error reading from terminal: %s", err))
		}

		if k.Special == terminal.KeyNoSpl {
			switch k.Regular {
			case '\r':
			//	session.Editor.InsertNewline()

			case ctrlKey('l'), '\x1b':
				break

			case ctrlKey('q'):
				//session.Quit()
				SafeExit(nil)

			case ctrlKey('s'):
				//session.Save()

			case ctrlKey('f'):
				//editorFind()

			case ctrlKey('h'), 127:
				//session.Editor.DelChar()

			default:
				//session.Editor.InsertChar(k.Regular)
			}
		} else {
			switch k.Special {

			case terminal.KeyArrowDown, terminal.KeyArrowLeft, terminal.KeyArrowRight, terminal.KeyArrowUp:
				//session.Editor.MoveCursor(k.Special)

			case terminal.KeyPageUp, terminal.KeyPageDown, terminal.KeyHome, terminal.KeyEnd:
				///session.Editor.MoveScreen(k.Special)

			case terminal.KeyDelete:
				//session.Editor.MoveCursor(terminal.KeyArrowRight)
				//session.Editor.DelChar()
			}
		}
	}
}

/*
func initEditor() error {
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

	b, err := rawmode.Enable()
	if err != nil {
		safeExit(err)
	}
	cfg.OrigTermCfg = b

	if err := initEditor(); err != nil {
		safeExit(err)
	}

	if len(os.Args) == 2 {
		rows, err := Open(os.Args[1])
		if err != nil {
			safeExit(err)
		}
		editor.Rows = rows
		editor.FileName = os.Args[1]
	}

	for {
		editorRefreshScreen()
		if err := editorProcessKeypress(); err != nil {
			safeExit(err)
		}
	}
}

func safeExit(err error) {
	fmt.Fprint(os.Stdout, "\x1b[2J")
	fmt.Fprint(os.Stdout, "\x1b[H")

	if err1 := rawmode.Restore(cfg.OrigTermCfg); err1 != nil {
		fmt.Fprintf(os.Stderr, "Error: disabling raw mode: %s\r\n", err)
	}

	if err == nil {
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Error: %s\r\n", err)
	os.Exit(1)
}
*/
