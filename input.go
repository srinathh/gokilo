package main

import (
	"gokilo/terminal"
)

func ctrlKey(b byte) rune {
	return rune(b & 0x1f)
}

func editorProcessKeypress() error {

	k, err := terminal.ReadKey()
	if err != nil {
		return err
	}

	switch k.Special {
	case terminal.KeyNoSpl:
		switch k.Regular {

		case ctrlKey('q'):
			/*
				if editor.Dirty && editor.QuitTimes > 0 {
					editorSetStatusMsg("WARNING!!! Unsaved changes. Press Ctrl-Q %d more times to quit.", editor.QuitTimes)
					editor.QuitTimes--
					return nil
				}*/
			SafeExit(nil)
		}

	}
	return nil
}

/*
func editorProcessKeypress() error {

	k, err := terminal.ReadKey()
	if err != nil {
		return err
	}

	if k.Special == terminal.KeyNoSpl {
		switch k.Regular {
		case '\r':
			editor.InsertNewline()
			break

		case ctrlKey('l'), '\x1b':
			break

		case ctrlKey('q'):
			if editor.Dirty && editor.QuitTimes > 0 {
				editorSetStatusMsg("WARNING!!! Unsaved changes. Press Ctrl-Q %d more times to quit.", editor.QuitTimes)
				editor.QuitTimes--
				return nil
			}
			safeExit(nil)
		case ctrlKey('s'):
			editor.Save()

		case ctrlKey('f'):
			editorFind()

		case ctrlKey('h'), 127:
			editor.DelChar()
		default:
			editor.InsertChar(k.Regular)
		}
	} else {
		switch k.Special {

		case terminal.KeyArrowDown, terminal.KeyArrowLeft, terminal.KeyArrowRight, terminal.KeyArrowUp:
			editorMoveCursor(k.Special)

		case terminal.KeyPageUp:
			editor.Cy = editor.RowOffset
			for j := 0; j < cfg.ScreenRows; j++ {
				editorMoveCursor(terminal.KeyArrowUp)
			}
		case terminal.KeyPageDown:
			editor.Cy = editor.RowOffset + cfg.ScreenRows - 1
			if editor.Cy > len(editor.Rows) {
				editor.Cy = len(editor.Rows)
			}
			for j := 0; j < cfg.ScreenRows; j++ {
				editorMoveCursor(terminal.KeyArrowDown)
			}
		case terminal.KeyHome:
			editor.Cx = 0
		case terminal.KeyEnd:
			if editor.Cy < len(editor.Rows) {
				editor.Cx = len(editor.Rows[editor.Cy])
			}

		case terminal.KeyDelete:
			editorMoveCursor(terminal.KeyArrowRight)
			editor.DelChar()
		}
	}

	editor.QuitTimes = kiloQuitTimes
	return nil
}

func editorMoveCursor(key int) {

	pastEOF := editor.Cy >= len(editor.Rows)

	switch key {
	case terminal.KeyArrowLeft:
		if editor.Cx > 0 {
			editor.Cx--
		} else if editor.Cy > 0 {
			editor.Cy--
			editor.Cx = len(editor.Rows[editor.Cy])
		}
	case terminal.KeyArrowRight:
		// right moves only if we're within a valid line.
		// for past EOF, there's no movement
		if !pastEOF {
			if editor.Cx < len(editor.Rows[editor.Cy]) {
				editor.Cx++
			} else if editor.Cx == len(editor.Rows[editor.Cy]) {
				editor.Cy++
				editor.Cx = 0
			}
		}
	case terminal.KeyArrowDown:
		if editor.Cy < len(editor.Rows) {
			editor.Cy++
		}
	case terminal.KeyArrowUp:
		if editor.Cy > 0 {
			editor.Cy--
		}
	}

	// we may have moved to a different row, so reset conditions
	pastEOF = editor.Cy >= len(editor.Rows)

	rowLen := 0
	if !pastEOF {
		rowLen = len(editor.Rows[editor.Cy])
	}

	if editor.Cx > rowLen {
		editor.Cx = rowLen
	}
}

type editorPromptCallback func(string, terminal.Key)

func editorPrompt(prompt string, callback editorPromptCallback) string {

	buf := ""

	for {
		editorSetStatusMsg(prompt, buf)
		editorRefreshScreen()

		k, err := terminal.ReadKey()
		if err != nil {
			return ""
		}

		switch {
		case k.Special == terminal.KeyDelete || k.Regular == 127 || k.Regular == ctrlKey('h'):
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
		case k.Regular == 27:
			editorSetStatusMsg("")
			if callback != nil {
				callback(buf, k)
			}
			return ""
		case k.Regular == '\r':
			if len(buf) != 0 {
				editorSetStatusMsg("")
				if callback != nil {
					callback(buf, k)
				}
				return buf
			}
		default:
			if k.Regular != ctrlKey('c') && k.Regular < 128 && k.Special == terminal.KeyNoSpl {
				buf = buf + string(k.Regular)
			}
			if callback != nil {
				callback(buf, k)
			}
		}
	}
}
*/
