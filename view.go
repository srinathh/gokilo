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

// Text expands tabs in an eRow to spaces
func Text(row ERow) ERow {
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
	txt := Text(row)
	if v.ColOffset > len(txt) {
		return ret
	}

	for j := v.ColOffset; j < len(txt) && j < v.ColOffset+v.ScreenCols; j++ {
		ret[j-v.ColOffset] = txt[j]
	}
	return ret
}

// RefreshScreen redraws the editing session on Screen
func (v *View) RefreshScreen(s *Session) {
	// clear screen
	//editorScroll()
	ab := bytes.Buffer{}

	// hide cursor
	fmt.Fprint(&ab, "\x1b[?25l")

	// move cursor to top left
	fmt.Fprint(&ab, "\x1b[H")

	v.DrawRows(&ab, s)
	//editorDrawStatusBar(&ab)
	//editorDrawStatusMsg(&ab)

	// reposition cursor
	//fmt.Fprintf(&ab, "\x1b[%d;%dH", editor.Cy-editor.RowOffset+1, editor.Rx-editor.ColOffset+1)

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

// DrawRows draws the editor rows on the screen
func (v *View) DrawRows(ab *bytes.Buffer, s *Session) {
	emptyRow := ERow("~")

	for y := 0; y < v.ScreenRows-1; y++ {
		fileRow := y + v.RowOffset
		if fileRow >= len(s.Editor.Rows) {
			fmt.Fprint(ab, string(v.ScreenText(emptyRow)))
		} else {
			fmt.Fprint(ab, string(v.ScreenText(s.Editor.Rows[fileRow]))) //editor.Rows[fileRow].ScreenText(editor.ColOffset, cfg.ScreenCols)))
		}
		fmt.Fprint(ab, "\r\n")
	}
}

/*
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
*/
