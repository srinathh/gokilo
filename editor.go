package main

import "fmt"

// ERow represents a line of text in a file
type ERow []rune

// Text expands tabs in an eRow to spaces
func (row ERow) Text() ERow {
	dest := []rune{}
	for _, r := range row {
		switch r {
		case '\t':
			dest = append(dest, tabSpaces...)
		default:
			dest = append(dest, r)
		}
	}
	return dest
}

// CxToRx transforms cursor positions to account for tab stops
func (row ERow) CxToRx(cx int) int {
	rx := 0
	for j := 0; j < cx; j++ {
		if row[j] == '\t' {
			rx = (rx + kiloTabStop - 1) - (rx % kiloTabStop)
		}
		rx++
	}
	return rx
}

// Editor represents the data in the being edited in memory
type Editor struct {
	Cx, Cy   int    // Cx and Cy represent current cursor position
	Rows     []ERow // Rows represent the textual data
	Dirty    bool   // has the file been edited
	FileName string // the path to the file being edited. Could be empty string
}

// NewEditor returns a new blank editor
func NewEditor() *Editor {
	return &Editor{
		FileName: "",
		Dirty:    false,
		Cx:       0,
		Cy:       0,
		Rows:     []ERow{},
	}
}

// NewEditorFromFile creates an editor from a file system file
func NewEditorFromFile(filename string) (*Editor, error) {

	rows := []ERow{}

	if filename != "" {
		var err error
		if rows, err = Open(filename); err != nil {
			return nil, fmt.Errorf("Error opening file %s: %v", filename, err)
		}
	}

	return &Editor{
		FileName: filename,
		Dirty:    false,
		Cx:       0,
		Cy:       0,
		Rows:     rows,
	}, nil
}

// CursorUp moves the cursor up as long as row number is non zero
func (e *Editor) CursorUp() {
	if e.Cy > 0 {
		e.Cy--
	}
	e.ResetX()
}

// CursorDown moves the cursor down till one line past the last line
func (e *Editor) CursorDown() {
	if e.Cy < len(e.Rows) {
		e.Cy++
	}
	e.ResetX()
}

// CursorLeft moves the cursor left. If at col 0 & any line other thant
// the first line, it moves to the previous line
func (e *Editor) CursorLeft() {

	if e.Cx > 0 {
		e.Cx--
	} else if e.Cy > 0 {
		e.Cy--
		e.Cx = len(e.Rows[e.Cy])
	}
}

// CursorRight moves the cursor right & wraps past EOL to col 0
func (e *Editor) CursorRight() {
	// right moves only if we're within a valid line.
	// for past EOF, there's no movement
	if e.Cy >= len(e.Rows) {
		return
	}
	if e.Cx < len(e.Rows[e.Cy]) {
		e.Cx++
	} else if e.Cx == len(e.Rows[e.Cy]) {
		e.Cy++
		e.Cx = 0
	}
}

// CursorEnd moves the cursor to end of line
func (e *Editor) CursorEnd() {
	if e.Cy < len(e.Rows) {
		e.Cx = len(e.Rows[e.Cy])
	}
}

// CursorHome moves the cursor to col 0
func (e *Editor) CursorHome() {
	e.Cx = 0
	/*
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
	*/
}

// ResetX sets the cursor X position to a valid position after moving y
func (e *Editor) ResetX() {

	// if we moved past last row, set cursor to 0
	if e.Cy >= len(e.Rows) {
		e.Cx = 0
		return
	}

	// we allow moving to 1 pos past the last character
	rowLen := len(e.Rows[e.Cy])
	if e.Cx > rowLen {
		e.Cx = len(e.Rows[e.Cy])
	}
}

/*
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
*/
