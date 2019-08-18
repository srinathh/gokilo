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
)

func editorProcessKeypress() error {

	b, err := editorReadKey()
	if err != nil {
		return err
	}

	switch b {
	case ctrlKey('q'):
		safeExit(nil)
	case keyArrowDown, keyArrowLeft, keyArrowRight, keyArrowUp:
		editorMoveCursor(b)
	case keyPageUp:
		for j := 0; j < cfg.screenRows; j++ {
			editorMoveCursor(keyArrowUp)
		}
	case keyPageDown:
		for j := 0; j < cfg.screenRows; j++ {
			editorMoveCursor(keyArrowDown)
		}
	case keyHome:
		cfg.cx = 0
	case keyEnd:
		cfg.cx = cfg.screenCols - 1
	}
	return nil
}

func editorMoveCursor(key int) {
	switch key {
	case keyArrowLeft:
		if cfg.cx > 0 {
			cfg.cx--
		}
	case keyArrowRight:
		if cfg.cx < cfg.screenCols-1 {
			cfg.cx++
		}
	case keyArrowDown:
		if cfg.cy < cfg.screenRows-1 {
			cfg.cy++
		}
	case keyArrowUp:
		if cfg.cy > 0 {
			cfg.cy--
		}
	}
}
