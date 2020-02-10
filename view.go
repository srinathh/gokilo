package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

var tabSpaces = []rune(strings.Repeat(" ", kiloTabStop))

const startMsg = "HELP: Ctrl + S to save | Ctrl + Q to exit | Ctrl + F to find"

const kiloTabStop = 4

const numStatusRows = 2

// View handles display of editor, status and prompts on screen
type View struct {
	ScreenRows    int
	ScreenCols    int
	RowOffset     int
	ColOffset     int
	StatusMsg     string // status message
	StatusMsgTime time.Time
	spaces        []rune
}

// NewView creates a view
func NewView(rows, cols int) *View {
	return &View{
		ScreenRows:    rows,
		ScreenCols:    cols,
		RowOffset:     0,
		ColOffset:     0,
		StatusMsg:     startMsg,
		StatusMsgTime: time.Now(),
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
	ret := spaces(v.ScreenCols)
	txt := row.Text()
	if v.ColOffset > len(txt) {
		return ret
	}

	for j := v.ColOffset; j < len(txt) && j < v.ColOffset+v.ScreenCols; j++ {
		ret[j-v.ColOffset] = txt[j]
	}
	return ret
}

// RefreshScreen redraws the editing session on Screen
func (v *View) RefreshScreen(e *Editor) {
	// clear screen
	rx := v.Scroll(e)
	ab := bytes.Buffer{}

	// hide cursor
	fmt.Fprint(&ab, "\x1b[?25l")

	// move cursor to top left
	fmt.Fprint(&ab, "\x1b[H")

	v.DrawRows(&ab, e)
	v.DrawStatusBar(&ab, e)
	//editorDrawStatusMsg(&ab)

	// reposition cursor
	fmt.Fprintf(&ab, "\x1b[%d;%dH", e.Cy-v.RowOffset+1, rx-v.ColOffset+1)

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

// Scroll scrolls the editor to capture the full view
func (v *View) Scroll(e *Editor) int {

	// if we're on the last line, cursor position should be 0
	rx := 0

	// find the screen x position after expanding tabs of the current row
	// as long as we're within the editor rows
	if e.Cy < len(e.Rows) {
		rx = e.Rows[e.Cy].CxToRx(e.Cx)
	}

	// if we have scrolled up beyond the current screen, move up
	if e.Cy < v.RowOffset {
		v.RowOffset = e.Cy
	}

	// if we have scrolled dwon below the screen, move down
	if e.Cy >= v.RowOffset+(v.ScreenRows-numStatusRows) {
		v.RowOffset = e.Cy - (v.ScreenRows - numStatusRows) + 1
	}

	// if we have scrolled left beyond hte screen, move our coloffset
	if rx < v.ColOffset {
		v.ColOffset = rx
	}

	if rx >= v.ColOffset+v.ScreenCols {
		v.ColOffset = rx - v.ScreenCols + 1
	}

	return rx
}

// DrawRows draws the editor rows on the screen
func (v *View) DrawRows(ab *bytes.Buffer, e *Editor) {
	emptyRow := ERow("~")

	for y := 0; y < v.ScreenRows-numStatusRows; y++ {
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
	rightStatusString := fmt.Sprintf("%dc %d/%dr", e.Cx+1, e.Cy+1, len(e.Rows))
	numSpaces := v.ScreenCols - len(leftStatusString) - len(rightStatusString)

	if numSpaces >= 0 {
		fmt.Fprint(ab, leftStatusString+strings.Repeat(" ", numSpaces)+rightStatusString)
	} else {
		fmt.Fprint(ab, (leftStatusString + rightStatusString)[:v.ScreenCols])
	}

	fmt.Fprint(ab, "\x1b[m")
	fmt.Fprint(ab, "\r\n")
}
