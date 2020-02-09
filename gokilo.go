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

	// parse config flags & parameters
	flag.Parse()
	filename := flag.Arg(0)

	// enable raw mode
	origCfg, err := rawmode.Enable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error enabling raw mode: %v", err)
		os.Exit(1)
	}

	SafeExit = func(error) { safeExit(origCfg, err) }

	// get the screen dimensions and create a view
	rows, cols, err := rawmode.GetWindowSize()
	if err != nil {
		SafeExit(fmt.Errorf("couldn't get window size: %v", err))
	}
	v := NewView(rows, cols)

	// create the editor
	var e *Editor
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
		// redraw the screen
		v.RefreshScreen(e)

		// read key
		k, err := terminal.ReadKey()
		if err != nil {
			SafeExit(fmt.Errorf("Error reading from terminal: %s", err))
		}

		// dispatch the key
		dispatchKey(k, nil, v, e)

	}
}

func dispatchKey(k terminal.Key, s *Session, v *View, e *Editor) {

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
