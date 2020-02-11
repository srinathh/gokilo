package main

import (
	"fmt"
)

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

// CursorPageUp moves the cursor one screen up
func (e *Editor) CursorPageUp(screenRows int, rowOffset int) {
	e.Cy = rowOffset
	for j := 0; j < screenRows; j++ {
		e.CursorUp()
	}
}

// CursorPageDown moves the cursor one screen down
func (e *Editor) CursorPageDown(screenRows int, rowOffset int) {
	e.Cy = rowOffset + screenRows - 1
	if e.Cy > len(e.Rows) {
		e.Cy = len(e.Rows)
	}
	for j := 0; j < screenRows; j++ {
		e.CursorDown()
	}
}

// CursorHome moves the cursor to col 0
func (e *Editor) CursorHome() {
	e.Cx = 0
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

// InsertChar inserts a character at a given location
func (e *Editor) InsertChar(c rune) {

	// if we're at the last line, insert a new row
	if e.Cy == len(e.Rows) {
		e.InsertRow(len(e.Rows), "")
	}

	// store a reference to the working row to improve readability
	src := e.Rows[e.Cy]

	dest := make([]rune, len(src)+1)
	copy(dest, src[:e.Cx])
	copy(dest[e.Cx+1:], src[e.Cx:])
	dest[e.Cx] = c

	e.Rows[e.Cy] = dest
	e.Dirty = true
	e.Cx++
}

// InsertRow inserts a row at a given index
func (e *Editor) InsertRow(rowidx int, s string) {
	if rowidx < 0 || rowidx > len(e.Rows) {
		return
	}

	row := []rune(s)

	e.Rows = append(e.Rows, ERow{})
	copy(e.Rows[rowidx+1:], e.Rows[rowidx:])
	e.Rows[rowidx] = row

	e.Dirty = true
}

// InsertNewline inserts a new line at the cursor position
func (e *Editor) InsertNewline() {
	if e.Cx == 0 {
		e.InsertRow(e.Cy, "")
		e.Dirty = true
		return
	}

	moveChars := string(e.Rows[e.Cy][e.Cx:])
	e.Rows[e.Cy] = e.Rows[e.Cy][:e.Cx]
	e.InsertRow(e.Cy+1, moveChars)
	e.Dirty = true

	e.Cy++
	e.Cx = 0
}
