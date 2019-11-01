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
	cfg.StatusMsg = fmt.Sprintf(format, a...)
	cfg.StatusMsgTime = time.Now()
}

func editorDrawStatusMsg(ab *bytes.Buffer) {
	fmt.Fprint(ab, "\x1b[K") // clear the line
	if time.Now().Sub(cfg.StatusMsgTime).Seconds() < 5 {
		if len(cfg.StatusMsg) < cfg.ScreenCols {
			fmt.Fprint(ab, cfg.StatusMsg)
		} else {
			fmt.Fprint(ab, cfg.StatusMsg[:cfg.ScreenCols])
		}
	}
}

func editorDrawRows(ab *bytes.Buffer) {
	for y := 0; y < cfg.ScreenRows; y++ {

		fileRow := y + editor.RowOffset

		if fileRow >= len(editor.Rows) {
			// print welcome message only if there is no file being edited
			if len(editor.Rows) == 0 && y == cfg.ScreenRows/3 {
				welcomeMsg := fmt.Sprintf("Kilo Editor -- version %s", kiloVersion)
				welcomeLen := len(welcomeMsg)

				// if the message is too long to fit, truncate
				if welcomeLen > cfg.ScreenCols {
					welcomeMsg = welcomeMsg[:cfg.ScreenCols]
					welcomeLen = cfg.ScreenCols
				}
				padding := (cfg.ScreenCols - welcomeLen) / 2

				// if there is at least 1 padding required, use the Tilde to start line
				if padding > 0 {
					fmt.Fprint(ab, "~")
					padding--
				}

				// add appropriate number of spaces
				for i := 0; i < padding; i++ {
					fmt.Fprint(ab, " ")
				}
				fmt.Fprint(ab, welcomeMsg)

			} else {
				fmt.Fprint(ab, "~")
			}
		} else {
			/*
				rowText := editor.Rows[fileRow].Text()
				rowSize := len(rowText) - editor.ColOffset
				if rowSize < 0 {
					rowSize = 0
				}
				if rowSize > cfg.ScreenCols {
					rowSize = cfg.ScreenCols
				}
				if rowSize > 0 {
					fmt.Fprint(ab, string(rowText[editor.ColOffset:editor.ColOffset+rowSize]))
				}
			*/
			fmt.Fprint(ab, string(editor.Rows[fileRow].ScreenText(editor.ColOffset, cfg.ScreenCols)))
		}

		// clear to end of line
		fmt.Fprint(ab, "\x1b[K")

		fmt.Fprint(ab, "\r\n")
	}
}
