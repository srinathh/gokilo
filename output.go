package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

func editorRefreshScreen() {
	// clear screen
	editorScroll()
	ab := bytes.Buffer{}

	// hide cursor
	fmt.Fprint(&ab, "\x1b[?25l")

	// move cursor to top left
	fmt.Fprint(&ab, "\x1b[H")

	editorDrawRows(&ab)
	editorDrawStatusBar(&ab)
	editorDrawStatusMsg(&ab)

	// reposition cursor
	fmt.Fprintf(&ab, "\x1b[%d;%dH", editor.Cy-editor.RowOffset+1, editor.Rx-editor.ColOffset+1)

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

func editorScroll() {
	editor.Rx = 0
	if editor.Cy < len(editor.Rows) {
		editor.Rx = editor.Rows[editor.Cy].CxToRx(editor.Cx)
	}

	if editor.Cy < editor.RowOffset {
		editor.RowOffset = editor.Cy
	}

	if editor.Cy >= editor.RowOffset+cfg.ScreenRows {
		editor.RowOffset = editor.Cy - cfg.ScreenRows + 1
	}

	if editor.Rx < editor.ColOffset {
		editor.ColOffset = editor.Rx
	}

	if editor.Rx >= editor.ColOffset+cfg.ScreenCols {
		editor.ColOffset = editor.Rx - cfg.ScreenCols + 1
	}
}

func editorDrawStatusBar(ab *bytes.Buffer) {
	fmt.Fprint(ab, "\x1b[7m")

	fileName := editor.FileName
	if fileName == "" {
		fileName = "No Name"
	}
	dirtyChar := ' '
	if editor.Dirty {
		dirtyChar = '*'
	}

	leftStatusString := fmt.Sprintf("%c%.20s - %d lines", dirtyChar, fileName, len(editor.Rows))
	rightStatusString := fmt.Sprintf("%dc %d/%dr", editor.Cx+1, editor.Cy+1, len(editor.Rows))
	numSpaces := cfg.ScreenCols - len(leftStatusString) - len(rightStatusString)

	if numSpaces >= 0 {
		fmt.Fprint(ab, leftStatusString+strings.Repeat(" ", numSpaces)+rightStatusString)
	} else {
		fmt.Fprint(ab, (leftStatusString + rightStatusString)[:cfg.ScreenCols])
	}

	fmt.Fprint(ab, "\x1b[m")
	fmt.Fprint(ab, "\r\n")
}

func editorSetStatusMsg(format string, a ...interface{}) {
	editor.StatusMsg = fmt.Sprintf(format, a...)
	editor.StatusMsgTime = time.Now()
}

func editorDrawStatusMsg(ab *bytes.Buffer) {
	fmt.Fprint(ab, "\x1b[K") // clear the line
	if time.Now().Sub(editor.StatusMsgTime).Seconds() < 5 {
		if len(editor.StatusMsg) < cfg.ScreenCols {
			fmt.Fprint(ab, editor.StatusMsg)
		} else {
			fmt.Fprint(ab, editor.StatusMsg[:cfg.ScreenCols])
		}
	}
}

func editorDrawRows(ab *bytes.Buffer) {
	emptyRow := ERow("~")

	for y := 0; y < cfg.ScreenRows; y++ {

		fileRow := y + editor.RowOffset

		if fileRow >= len(editor.Rows) {
			// print welcome message only if there is no file being edited
			if len(editor.Rows) == 0 && y == cfg.ScreenRows/3 {
				welcomeMsg := ERow(fmt.Sprintf("~ Kilo Editor -- version %s", kiloVersion))
				fmt.Fprint(ab, string(welcomeMsg.ScreenText(0, cfg.ScreenCols)))
			} else {
				fmt.Fprint(ab, string(emptyRow.ScreenText(0, cfg.ScreenCols)))
			}
		} else {
			fmt.Fprint(ab, string(editor.Rows[fileRow].ScreenText(editor.ColOffset, cfg.ScreenCols)))
		}

		fmt.Fprint(ab, "\r\n")
	}
}
