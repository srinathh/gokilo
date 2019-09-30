package main

const (
	keyArrowUp    = 1000
	keyArrowDown  = 1001
	keyArrowLeft  = 1002
	keyArrowRight = 1003
	keyPageUp     = 1004
	keyPageDown   = 1005
	keyHome       = 1006
	keyEnd        = 1007
	keyDelete     = 1008
	keyBackSpace  = 127
)

func editorProcessKeypress() error {

	b := editorReadKey()

	switch b {
	case '\r':
		editorInsertNewline()
		break

	case ctrlKey('l'), '\x1b':
		break

	case ctrlKey('q'):
		if cfg.dirty && cfg.quitTimes > 0 {
			editorSetStatusMsg("WARNING!!! Unsaved changes. Press Ctrl-Q %d more times to quit.", cfg.quitTimes)
			cfg.quitTimes--
			return nil
		}
		safeExit(nil)
	case keyArrowDown, keyArrowLeft, keyArrowRight, keyArrowUp:
		editorMoveCursor(b)
	case keyPageUp:
		cfg.cy = cfg.rowOffset
		for j := 0; j < cfg.screenRows; j++ {
			editorMoveCursor(keyArrowUp)
		}
	case keyPageDown:
		cfg.cy = cfg.rowOffset + cfg.screenRows - 1
		if cfg.cy > len(cfg.rows) {
			cfg.cy = len(cfg.rows)
		}
		for j := 0; j < cfg.screenRows; j++ {
			editorMoveCursor(keyArrowDown)
		}
	case keyHome:
		cfg.cx = 0

	case ctrlKey('s'):
		editorSave()

	case ctrlKey('f'):
		editorFind()

	case keyEnd:
		if cfg.cy < len(cfg.rows) {
			cfg.cx = len(cfg.rows[cfg.cy])
		}

	case keyBackSpace, ctrlKey('h'):
		editorDelChar()

	case keyDelete:
		editorMoveCursor(keyArrowRight)
		editorDelChar()

	default:
		editorInsertChar(b)
	}
	cfg.quitTimes = kiloQuitTimes
	return nil
}

func editorMoveCursor(key int) {

	pastEOF := cfg.cy >= len(cfg.rows)

	switch key {
	case keyArrowLeft:
		if cfg.cx > 0 {
			cfg.cx--
		} else if cfg.cy > 0 {
			cfg.cy--
			cfg.cx = len(cfg.rows[cfg.cy])
		}
	case keyArrowRight:
		// right moves only if we're within a valid line.
		// for past EOF, there's no movement
		if !pastEOF {
			if cfg.cx < len(cfg.rows[cfg.cy]) {
				cfg.cx++
			} else if cfg.cx == len(cfg.rows[cfg.cy]) {
				cfg.cy++
				cfg.cx = 0
			}
		}
	case keyArrowDown:
		if cfg.cy < len(cfg.rows) {
			cfg.cy++
		}
	case keyArrowUp:
		if cfg.cy > 0 {
			cfg.cy--
		}
	}

	// we may have moved to a different row, so reset conditions
	pastEOF = cfg.cy >= len(cfg.rows)

	rowLen := 0
	if !pastEOF {
		rowLen = len(cfg.rows[cfg.cy])
	}

	if cfg.cx > rowLen {
		cfg.cx = rowLen
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
