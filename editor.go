package main

// Editor represents a file being edited
type Editor struct {
	Cx, Cy    int
	Rx        int
	Rows      []ERow
	RowOffset int
	ColOffset int
	FileName  string
	Dirty     bool
	QuitTimes int
	LastMatch int
	Direction int
}

// DelChar deletes a character at current cursor location
func (e *Editor) DelChar() {
	// if the cursor is in the empty line at the end, do nothing (why?)
	if e.Cy == len(e.Rows) {
		return
	}

	// if at the beginning of the text, then do nothing
	if e.Cx == 0 && e.Cy == 0 {
		return
	}

	// different handling for at the beginning of the line or middle of line
	if e.Cx > 0 {
		e.Rows[e.Cy] = e.Rows[e.Cy].DelChar(e.Cx - 1)
		e.Cx--
	} else {
		e.Cx = len(e.Rows[e.Cy-1])
		e.Rows[e.Cy-1] = append(e.Rows[e.Cy-1], e.Rows[e.Cy]...)
		e.DelRow(e.Cy)
		e.Cy--
	}
	e.Dirty = true
}

// DelRow deletes a given row
func (e *Editor) DelRow(rowidx int) {
	if rowidx < 0 || rowidx >= len(e.Rows) {
		return
	}

	copy(e.Rows[rowidx:], e.Rows[rowidx+1:])
	e.Rows = e.Rows[:len(e.Rows)-1]
	e.Dirty = true
}

// InsertChar inserts a character at a given location
func (e *Editor) InsertChar(c rune) {
	if e.Cy == len(e.Rows) {
		e.InsertRow(len(e.Rows), "")
	}
	e.Rows[e.Cy] = e.Rows[e.Cy].InsertChar(e.Cx, c)
	e.Dirty = true
	e.Cx++
}

// InsertNewline inserts a new line at the cursor position
func (e *Editor) InsertNewline() {
	if e.Cx == 0 {
		e.InsertRow(e.Cy, "")
		return
	}

	moveChars := string(e.Rows[e.Cy][e.Cx:])

	e.Rows[e.Cy] = e.Rows[e.Cy][:e.Cx]

	e.InsertRow(e.Cy+1, moveChars)

	e.Cy++
	e.Cx = 0
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

// Save gets a filename if required and saves a file
func (e *Editor) Save() {

	if e.FileName == "" {
		e.FileName = editorPrompt("Save as: %s", nil)
		if editor.FileName == "" {
			editorSetStatusMsg("Save aborted!")
			return
		}
	}

	if err := Save(e.Rows, e.FileName); err != nil {
		editorSetStatusMsg("ERROR SAVING: %s", err)
	} else {
		editorSetStatusMsg("SAVED FILE: %s", e.FileName)
		e.Dirty = false
	}
}
