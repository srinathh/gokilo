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
		FileName: "",
		Dirty:    false,
		Cx:       0,
		Cy:       0,
		Rows:     rows,
	}, nil
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
*/
