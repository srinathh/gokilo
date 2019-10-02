package main

func editorProcessKeypress() error {

	b := editorReadKey()

	switch b {
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
	case keyArrowDown, keyArrowLeft, keyArrowRight, keyArrowUp:
		editorMoveCursor(b)
	case keyPageUp:
		editor.Cy = editor.RowOffset
		for j := 0; j < cfg.ScreenRows; j++ {
			editorMoveCursor(keyArrowUp)
		}
	case keyPageDown:
		editor.Cy = editor.RowOffset + cfg.ScreenRows - 1
		if editor.Cy > len(editor.Rows) {
			editor.Cy = len(editor.Rows)
		}
		for j := 0; j < cfg.ScreenRows; j++ {
			editorMoveCursor(keyArrowDown)
		}
	case keyHome:
		editor.Cx = 0

	case ctrlKey('s'):
		editor.Save()

	case ctrlKey('f'):
		editorFind()

	case keyEnd:
		if editor.Cy < len(editor.Rows) {
			editor.Cx = len(editor.Rows[editor.Cy])
		}

	case keyBackSpace, ctrlKey('h'):
		editor.DelChar()

	case keyDelete:
		editorMoveCursor(keyArrowRight)
		editor.DelChar()

	default:
		editor.InsertChar(b)
	}
	editor.QuitTimes = kiloQuitTimes
	return nil
}

func editorMoveCursor(key int) {

	pastEOF := editor.Cy >= len(editor.Rows)

	switch key {
	case keyArrowLeft:
		if editor.Cx > 0 {
			editor.Cx--
		} else if editor.Cy > 0 {
			editor.Cy--
			editor.Cx = len(editor.Rows[editor.Cy])
		}
	case keyArrowRight:
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
	case keyArrowDown:
		if editor.Cy < len(editor.Rows) {
			editor.Cy++
		}
	case keyArrowUp:
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

type editorPromptCallback func(string, int)

func editorPrompt(prompt string, callback editorPromptCallback) string {

	buf := ""

	for {
		editorSetStatusMsg(prompt, buf)
		editorRefreshScreen()

		c := editorReadKey()
		switch c {
		case keyDelete, keyBackSpace, ctrlKey('h'):
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
			}
		case '\x1b':
			editorSetStatusMsg("")
			if callback != nil {
				callback(buf, c)
			}
			return ""
		case '\r':
			if len(buf) != 0 {
				editorSetStatusMsg("")
				if callback != nil {
					callback(buf, c)
				}
				return buf
			}
		default:
			if c != ctrlKey('c') && c < 128 {
				buf = buf + string(c)
			}
			if callback != nil {
				callback(buf, c)
			}
		}
	}
}
