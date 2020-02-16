package main

import (
	"gokilo/terminal"
)

// Session stores state pertaining to the editing session
type Session struct {
	Prompt      *LineEditor
	Editor      *Editor
	View        *View
	StatusMsg   string
	OrigTermCfg []byte
	State       editorState
}

// startQuitPrompt will set the status message
func startQuitPrompt() {
	s.State = stateQuitPrompt
	s.StatusMsg = "Unsaved Chages! Press Ctrl+Q again to quit, any other key to cancel"
	s.Prompt = NewLineEditor()
}

// startQuitPrompt will set the status message
func endQuitPromt() {
	s.State = stateEditing
	s.StatusMsg = ""
	s.Prompt = nil
}

func editingStateDispatch(k terminal.Key, v *View, e *Editor) {

	if k.Special == terminal.KeyNoSpl {
		switch k.Regular {
		case '\r':
			e.InsertNewline()
			break

		case ctrlKey('h'), 127:
			e.DelChar()

		default:
			e.InsertChar(k.Regular)
		}
	} else {
		switch k.Special {

		case terminal.KeyArrowDown:
			e.CursorDown()

		case terminal.KeyArrowLeft:
			e.CursorLeft()

		case terminal.KeyArrowRight:
			e.CursorRight()

		case terminal.KeyArrowUp:
			e.CursorUp()

		case terminal.KeyHome:
			e.CursorHome()

		case terminal.KeyEnd:
			e.CursorEnd()

		case terminal.KeyPageUp:
			e.CursorPageUp(v.ScreenRows, v.RowOffset)

		case terminal.KeyPageDown:
			e.CursorPageDown(v.ScreenRows, v.RowOffset)

		case terminal.KeyDelete:
			e.CursorRight()
			e.DelChar()
		}
	}
}
