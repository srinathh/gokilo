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
	fmt.Fprintf(&ab, "\x1b[%d;%dH", cfg.cy-cfg.rowOffset+1, cfg.rx-cfg.colOffset+1)

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

func editorScroll() {
	cfg.rx = 0
	if cfg.cy < len(cfg.rows) {
		cfg.rx = cfg.rows[cfg.cy].editorRowCxToRx(cfg.cx)
	}

	if cfg.cy < cfg.rowOffset {
		cfg.rowOffset = cfg.cy
	}

	if cfg.cy >= cfg.rowOffset+cfg.screenRows {
		cfg.rowOffset = cfg.cy - cfg.screenRows + 1
	}

	if cfg.rx < cfg.colOffset {
		cfg.colOffset = cfg.rx
	}

	if cfg.rx >= cfg.colOffset+cfg.screenCols {
		cfg.colOffset = cfg.rx - cfg.screenCols + 1
	}
}

func editorDrawStatusBar(ab *bytes.Buffer) {
	fmt.Fprint(ab, "\x1b[7m")

	fileName := cfg.fileName
	if fileName == "" {
		fileName = "No Name"
	}
	dirtyChar := ' '
	if cfg.dirty {
		dirtyChar = '*'
	}

	leftStatusString := fmt.Sprintf("%c%.20s - %d lines", dirtyChar, fileName, len(cfg.rows))
	rightStatusString := fmt.Sprintf("%dc %d/%dr", cfg.cx+1, cfg.cy+1, len(cfg.rows))
	numSpaces := cfg.screenCols - len(leftStatusString) - len(rightStatusString)

	if numSpaces >= 0 {
		fmt.Fprint(ab, leftStatusString+strings.Repeat(" ", numSpaces)+rightStatusString)
	} else {
		fmt.Fprint(ab, (leftStatusString + rightStatusString)[:cfg.screenCols])
	}

	fmt.Fprint(ab, "\x1b[m")
	fmt.Fprint(ab, "\r\n")
}

func editorSetStatusMsg(format string, a ...interface{}) {
	cfg.statusMsg = fmt.Sprintf(format, a...)
	cfg.statusMsgTime = time.Now()
}

func editorDrawStatusMsg(ab *bytes.Buffer) {
	fmt.Fprint(ab, "\x1b[K") // clear the line
	if time.Now().Sub(cfg.statusMsgTime).Seconds() < 5 {
		if len(cfg.statusMsg) < cfg.screenCols {
			fmt.Fprint(ab, cfg.statusMsg)
		} else {
			fmt.Fprint(ab, cfg.statusMsg[:cfg.screenCols])
		}
	}
}

func editorDrawRows(ab *bytes.Buffer) {
	for y := 0; y < cfg.screenRows; y++ {

		fileRow := y + cfg.rowOffset

		if fileRow >= len(cfg.rows) {
			// print welcome message only if there is no file being edited
			if len(cfg.rows) == 0 && y == cfg.screenRows/3 {
				welcomeMsg := fmt.Sprintf("Kilo Editor -- version %s", kiloVersion)
				welcomeLen := len(welcomeMsg)

				// if the message is too long to fit, truncate
				if welcomeLen > cfg.screenCols {
					welcomeMsg = welcomeMsg[:cfg.screenCols]
					welcomeLen = cfg.screenCols
				}
				padding := (cfg.screenCols - welcomeLen) / 2

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
			rowText := cfg.rows[fileRow].Text()
			rowSize := len(rowText) - cfg.colOffset
			if rowSize < 0 {
				rowSize = 0
			}
			if rowSize > cfg.screenCols {
				rowSize = cfg.screenCols
			}
			if rowSize > 0 {
				fmt.Fprint(ab, string(rowText[cfg.colOffset:cfg.colOffset+rowSize]))
			}
		}

		// clear to end of line
		fmt.Fprint(ab, "\x1b[K")

		fmt.Fprint(ab, "\r\n")
	}
}
