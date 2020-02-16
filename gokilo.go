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

type editorState int

const (
	stateEditing editorState = iota
	stateSavePrompt
	stateQuitPrompt
	stateFindPrompt
	stateFindNav
)

const kiloVersion = "0.0.2"

// SafeExit restores terminal using the original terminal config stored
// in the global session variable
func SafeExit(err error) {
	fmt.Fprint(os.Stdout, "\x1b[2J\x1b[H")

	if err1 := rawmode.Restore(s.OrigTermCfg); err1 != nil {
		fmt.Fprintf(os.Stderr, "Error: disabling raw mode: %s\r\n", err)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\r\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

var s = Session{}

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
	s.OrigTermCfg = origCfg

	// get the screen dimensions and create a view
	rows, cols, err := rawmode.GetWindowSize()
	if err != nil {
		SafeExit(fmt.Errorf("couldn't get window size: %v", err))
	}
	s.View = NewView(rows, cols)

	// create the editor
	if flag.Arg(0) == "" {
		s.Editor = NewEditor()
	} else {
		s.Editor, err = NewEditorFromFile(filename)
		if err != nil {
			SafeExit(fmt.Errorf("couldn't open file %s: %v", filename, err))
		}
	}

	s.StatusMsg = startMsg
	s.State = stateEditing

	for {
		s.View.RefreshScreen(s.Editor, s.StatusMsg, s.Prompt)

		// read key
		k, err := terminal.ReadKey()
		if err != nil {
			SafeExit(fmt.Errorf("Error reading from terminal: %s", err))
		}

		switch s.State {
		case stateEditing:
			switch {
			case k.Regular == ctrlKey('Q'):
				if !s.Editor.Dirty {
					SafeExit(nil)
				}
				startQuitPrompt()
			default:
				editingStateDispatch(k, s.View, s.Editor)
			}
		case stateQuitPrompt:
			if k.Regular == ctrlKey('Q') {
				SafeExit(nil)
			}
			endQuitPromt()
		}

	}
}

/*

					stateEditing	stateSavePrompt		stateQuitPrompt		stateFindPromp		stateFindNav
stateEditing		any other key	Ctrl+S & NoFname	Ctrl+Q & Dirty		Ctrl+F
stateSavePrompt		Esc or Enter	any other key
stateQuitPrompt		any other key	Ctrl+Q
stateFindPrompt		Esc										any other key		Enter
stateFindNav		Esc or Enter


*/
