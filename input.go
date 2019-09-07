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

	b, err := editorReadKey()
	if err != nil {
		return err
	}

	switch b {
	case '\r':
		//tk
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

	case keyEnd:
		if cfg.cy < len(cfg.rows) {
			cfg.cx = len(cfg.rows[cfg.cy].chars)
		}

	case keyBackSpace, ctrlKey('h'):
		editorDelChar()

	case keyDelete:
		editorMoveCursor(keyArrowRight)
		editorDelChar()

	case ctrlKey('l'), '\x1b':
		editorSetStatusMsg("enter pressed")

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
			cfg.cx = len(cfg.rows[cfg.cy].chars)
		}
	case keyArrowRight:
		// right moves only if we're within a valid line.
		// for past EOF, there's no movement
		if !pastEOF {
			if cfg.cx < len(cfg.rows[cfg.cy].chars) {
				cfg.cx++
			} else if cfg.cx == len(cfg.rows[cfg.cy].chars) {
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
		rowLen = len(cfg.rows[cfg.cy].chars)
	}

	if cfg.cx > rowLen {
		cfg.cx = rowLen
	}
}
