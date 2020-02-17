package main

import (
	"fmt"
	"time"

	"github.com/srinathh/gokilo/terminal"
)

// Point is a utility struct
type Point struct {
	Col, Row int
}

// Session stores state pertaining to the editing session
type Session struct {
	Prompt            *LineEditor
	Editor            *Editor
	View              *View
	StatusMessage     string
	StatusMessageTime time.Time
	OrigTermCfg       []byte
	State             editorState
	FindPoints        []Point
	CurFindPoint      int
	BeforFindPoint    Point
}

func (s *Session) setStatusMessage(msg string) {
	s.StatusMessage = msg
	s.StatusMessageTime = time.Now()
}

func (s *Session) saveFile() {
	if err := Save(s.Editor.Rows, s.Editor.FileName); err != nil {
		s.setStatusMessage(fmt.Sprintf("Error saving %s: %s", s.Editor.FileName, err))
	} else {
		s.setStatusMessage(fmt.Sprintf("Saved %s", s.Editor.FileName))
		s.Editor.Dirty = false
	}
}

func (s *Session) startSavePrompt() {
	s.State = stateSavePrompt
	s.setStatusMessage("Enter filename: ")
	s.Prompt = NewLineEditor()
}

func (s *Session) endSavePrompt(save bool) {
	if save {
		s.Editor.FileName = string(s.Prompt.Row.Text())
		s.saveFile()
	} else {
		s.setStatusMessage("")
	}
	s.Prompt = nil
	s.State = stateEditing
}

func (s *Session) startFindPrompt() {
	s.State = stateFindPrompt
	s.setStatusMessage("Find: ")
	s.Prompt = NewLineEditor()
}

func (s *Session) endFindPrompt(findnav bool) {
	if !findnav {
		s.Prompt = nil
		s.setStatusMessage("")
		s.State = stateEditing
		return
	}
	s.startFindNav()
}

func (s *Session) startFindNav() {
	findText := s.Prompt.Row.Text()
	s.State = stateFindNav
	s.Prompt = nil

	s.FindPoints = s.Editor.Find(findText)
	if len(s.FindPoints) == 0 {
		s.setStatusMessage("No match found!")
		s.FindPoints = nil
		s.CurFindPoint = -1
		s.Prompt = nil
		s.State = stateEditing
		return
	}
	s.BeforFindPoint = s.Editor.Cursor
	s.CurFindPoint = 0
	s.Editor.SetCursor(s.FindPoints[s.CurFindPoint])
	s.setStatusMessage("Use arrow keys to move, ESC or ENTER to exit")
}

func (s *Session) endFindNav() {
	s.BeforFindPoint = Point{}
	s.FindPoints = nil
	s.CurFindPoint = -1
	s.Prompt = nil
	s.State = stateEditing
	s.setStatusMessage("")
}

/*Dispatch decides what action to take given the user's key and current state

					stateEditing	stateSavePrompt		stateQuitPrompt		stateFindPromp		stateFindNav
stateEditing		any other key	Ctrl+S & NoFname	Ctrl+Q & Dirty		Ctrl+F
stateSavePrompt		Esc or Enter	any other key
stateQuitPrompt		any other key	Ctrl+Q
stateFindPrompt		Esc										any other key		Enter
stateFindNav		Esc or Enter																arrow keys

*/
func (s *Session) Dispatch(k terminal.Key) {
	switch s.State {
	case stateEditing:
		switch {
		case k.Regular == ctrlKey('Q'):
			if !s.Editor.Dirty {
				SafeExit(nil)
			}
			s.startQuitPrompt()
		case k.Regular == ctrlKey('S'):
			if s.Editor.FileName != "" {
				s.saveFile()
			} else {
				s.startSavePrompt()
			}
		case k.Regular == ctrlKey('F'):
			s.startFindPrompt()
		default:
			s.editorDispatch(k)
		}

	case stateQuitPrompt:
		if k.Regular == ctrlKey('Q') {
			SafeExit(nil)
		}
		s.endQuitPromt()

	case stateSavePrompt:
		switch {
		case k.Regular == '\r':
			s.endSavePrompt(true)
		case k.Regular == 27:
			s.endSavePrompt(false)
		default:
			s.lineEditorDispatch(k)
		}

	case stateFindPrompt:
		switch {
		case k.Regular == '\r':
			s.endFindPrompt(true)
		case k.Regular == 27:
			s.endFindPrompt(false)
		default:
			s.lineEditorDispatch(k)
		}

	case stateFindNav:
		switch {
		case k.Special == terminal.KeyArrowUp, k.Special == terminal.KeyArrowLeft:
			s.CurFindPoint--
			if s.CurFindPoint < 0 {
				s.CurFindPoint = len(s.FindPoints) - 1
			}
			s.Editor.SetCursor(s.FindPoints[s.CurFindPoint])

		case k.Special == terminal.KeyArrowDown, k.Special == terminal.KeyArrowRight:
			s.CurFindPoint++
			if s.CurFindPoint >= len(s.FindPoints) {
				s.CurFindPoint = 0
			}
			s.Editor.SetCursor(s.FindPoints[s.CurFindPoint])
		case k.Regular == 27:
			s.Editor.SetCursor(s.BeforFindPoint)
			s.endFindNav()
		case k.Regular == '\r':
			s.endFindNav()
		}
	}
}

// startQuitPrompt will set the status message
func (s *Session) startQuitPrompt() {
	s.State = stateQuitPrompt
	s.setStatusMessage("Unsaved Chages! Press Ctrl+Q again to quit, any other key to cancel")
	s.Prompt = NewLineEditor()
}

// startQuitPrompt will set the status message
func (s *Session) endQuitPromt() {
	s.State = stateEditing
	s.setStatusMessage("")
	s.Prompt = nil
}

func (s *Session) editorDispatch(k terminal.Key) {

	if k.Special == terminal.KeyNoSpl {
		switch k.Regular {
		case '\r':
			s.Editor.InsertNewline()
			break

		case ctrlKey('h'), 127:
			s.Editor.DelChar()

		default:
			if k.Regular >= 32 {
				s.Editor.InsertChar(k.Regular)
			}
		}
	} else {
		switch k.Special {

		case terminal.KeyArrowDown:
			s.Editor.CursorDown()

		case terminal.KeyArrowLeft:
			s.Editor.CursorLeft()

		case terminal.KeyArrowRight:
			s.Editor.CursorRight()

		case terminal.KeyArrowUp:
			s.Editor.CursorUp()

		case terminal.KeyHome:
			s.Editor.CursorHome()

		case terminal.KeyEnd:
			s.Editor.CursorEnd()

		case terminal.KeyPageUp:
			s.Editor.CursorPageUp(s.View.ScreenSize.Row, s.View.RowOffset)

		case terminal.KeyPageDown:
			s.Editor.CursorPageDown(s.View.ScreenSize.Row, s.View.RowOffset)

		case terminal.KeyDelete:
			s.Editor.CursorRight()
			s.Editor.DelChar()
		}
	}
}

func (s *Session) lineEditorDispatch(k terminal.Key) {

	if k.Special == terminal.KeyNoSpl {
		switch k.Regular {
		case ctrlKey('h'), 127:
			s.Prompt.DelChar()

		default:
			if k.Regular >= 32 {
				s.Prompt.InsertChar(k.Regular)
			}
		}
	} else {
		switch k.Special {

		case terminal.KeyArrowLeft:
			s.Prompt.CursorLeft()

		case terminal.KeyArrowRight:
			s.Prompt.CursorRight()

		case terminal.KeyHome:
			s.Prompt.CursorHome()

		case terminal.KeyEnd:
			s.Prompt.CursorEnd()

		case terminal.KeyDelete:
			s.Prompt.CursorRight()
			s.Prompt.DelChar()
		}
	}
}
