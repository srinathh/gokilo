package main

// LineEditor represents the a line being edited
type LineEditor struct {
	Cx  int  // Cx and Cy represent current cursor position
	Row ERow // Rows represent the textual data
}

// NewLineEditor returns a new blank editor
func NewLineEditor() *LineEditor {
	return &LineEditor{
		Cx:  0,
		Row: ERow{},
	}
}

// CursorLeft moves the cursor left. If at col 0 & any line other thant
// the first line, it moves to the previous line
func (e *LineEditor) CursorLeft() {

	if e.Cx > 0 {
		e.Cx--
	}
}

// CursorRight moves the cursor right & wraps past EOL to col 0
func (e *LineEditor) CursorRight() {
	// right moves only if we're within a valid line.
	// for past EOF, there's no movement
	if e.Cx < len(e.Row) {
		e.Cx++
	}
}

// CursorEnd moves the cursor to end of line
func (e *LineEditor) CursorEnd() {
	e.Cx = len(e.Row)
}

// CursorHome moves the cursor to col 0
func (e *LineEditor) CursorHome() {
	e.Cx = 0
}

// InsertChar inserts a character at a given location
func (e *LineEditor) InsertChar(c rune) {

	// store a reference to the working row to improve readability
	src := e.Row

	dest := make([]rune, len(src)+1)
	copy(dest, src[:e.Cx])
	copy(dest[e.Cx+1:], src[e.Cx:])
	dest[e.Cx] = c

	e.Row = dest
	e.Cx++
}

// DelChar deletes a character at current cursor location
func (e *LineEditor) DelChar() {

	// different handling for at the beginning of the line or middle of line
	if e.Cx > 0 {
		row := e.Row
		copy(row[e.Cx-1:], row[e.Cx:])
		row = row[:len(row)-1]
		e.Row = row
		e.Cx--
	}
}
