package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

var tabSpaces = []rune(strings.Repeat(" ", kiloTabStop))

const startMsg = "HELP: Ctrl + S to save | Ctrl + Q to exit | Ctrl + F to find"

const kiloTabStop = 4

const numStatusRows = 2

// View handles display of editor, status and prompts on screen
type View struct {
	//ScreenRows int
	//ScreenCols int
	ScreenSize Point
	RowOffset  int
	ColOffset  int
	spaces     []rune
}

// NewView creates a view
func NewView(rows, cols int) *View {
	return &View{
		ScreenSize: Point{Row: rows, Col: cols},
		RowOffset:  0,
		ColOffset:  0,
	}
}

func spaces(width int) ERow {
	ret := make([]rune, width)
	for j := 0; j < width; j++ {
		ret[j] = ' '
	}
	return ret
}

// ScreenText returns a line of text offset to match screen window
func (v *View) ScreenText(row ERow) []rune {
	ret := spaces(v.ScreenSize.Col)
	txt := row.Text()
	if v.ColOffset > len(txt) {
		return ret
	}

	for j := v.ColOffset; j < len(txt) && j < v.ColOffset+v.ScreenSize.Col; j++ {
		ret[j-v.ColOffset] = txt[j]
	}
	return ret
}

// RefreshScreen redraws the editing session on Screen
func (v *View) RefreshScreen(e *Editor, statusMsg string, prompt *LineEditor) {
	// clear screen
	rx := v.Scroll(e)
	ab := bytes.Buffer{}

	// hide cursor
	fmt.Fprint(&ab, "\x1b[?25l")

	// move cursor to top left
	fmt.Fprint(&ab, "\x1b[H")

	v.DrawRows(&ab, e)
	v.DrawStatusBar(&ab, e)
	v.DrawStatusMsg(&ab, statusMsg)
	if prompt != nil {
		fmt.Fprint(&ab, string(prompt.Row.Text()))
		fmt.Fprintf(&ab, "\x1b[%d;%dH", v.ScreenSize.Row, len(statusMsg)+prompt.Cx+1)
	} else {
		fmt.Fprintf(&ab, "\x1b[%d;%dH", e.Cursor.Row-v.RowOffset+1, rx-v.ColOffset+1)
	}

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

// DrawStatusMsg draws the status message on screen
func (v *View) DrawStatusMsg(ab *bytes.Buffer, statusMsg string) {
	fmt.Fprint(ab, "\x1b[K") // clear the line
	if len(statusMsg) < v.ScreenSize.Col {
		fmt.Fprint(ab, statusMsg)
	} else {
		fmt.Fprint(ab, statusMsg[:v.ScreenSize.Col])
	}
}

// Scroll scrolls the editor to capture the full view
func (v *View) Scroll(e *Editor) int {

	// if we're on the last line, cursor position should be 0
	rx := 0

	// find the screen x position after expanding tabs of the current row
	// as long as we're within the editor rows
	if e.Cursor.Row < len(e.Rows) {
		rx = e.Rows[e.Cursor.Row].CxToRx(e.Cursor.Col)
	}

	// if we have scrolled up beyond the current screen, move up
	if e.Cursor.Row < v.RowOffset {
		v.RowOffset = e.Cursor.Row
	}

	// if we have scrolled dwon below the screen, move down
	if e.Cursor.Row >= v.RowOffset+(v.ScreenSize.Row-numStatusRows) {
		v.RowOffset = e.Cursor.Row - (v.ScreenSize.Row - numStatusRows) + 1
	}

	// if we have scrolled left beyond hte screen, move our coloffset
	if rx < v.ColOffset {
		v.ColOffset = rx
	}

	if rx >= v.ColOffset+v.ScreenSize.Col {
		v.ColOffset = rx - v.ScreenSize.Col + 1
	}

	return rx
}

// DrawRows draws the editor rows on the screen
func (v *View) DrawRows(ab *bytes.Buffer, e *Editor) {
	emptyRow := ERow("~")

	for y := 0; y < v.ScreenSize.Row-numStatusRows; y++ {
		fileRow := y + v.RowOffset
		if fileRow >= len(e.Rows) {
			fmt.Fprint(ab, string(v.ScreenText(emptyRow)))
		} else {
			fmt.Fprint(ab, string(v.ScreenText(e.Rows[fileRow]))) //editor.Rows[fileRow].ScreenText(editor.ColOffset, cfg.ScreenCols)))
		}
		fmt.Fprint(ab, "\r\n")
	}
}

// DrawStatusBar draws the status bar
func (v *View) DrawStatusBar(ab *bytes.Buffer, e *Editor) {
	fmt.Fprint(ab, "\x1b[7m")

	fileName := e.FileName
	if fileName == "" {
		fileName = "No Name"
	}
	dirtyChar := ' '
	if e.Dirty {
		dirtyChar = '*'
	}

	leftStatusString := fmt.Sprintf("%c%.20s - %d lines", dirtyChar, fileName, len(e.Rows))
	rightStatusString := fmt.Sprintf("%dc %d/%dr", e.Cursor.Col+1, e.Cursor.Row+1, len(e.Rows))
	numSpaces := v.ScreenSize.Col - len(leftStatusString) - len(rightStatusString)

	if numSpaces >= 0 {
		fmt.Fprint(ab, leftStatusString+strings.Repeat(" ", numSpaces)+rightStatusString)
	} else {
		fmt.Fprint(ab, (leftStatusString + rightStatusString)[:v.ScreenSize.Col])
	}

	fmt.Fprint(ab, "\x1b[m")
	fmt.Fprint(ab, "\r\n")
}
