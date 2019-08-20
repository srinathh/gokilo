package main

import (
	"bytes"
	"fmt"
	"os"
)

func editorRefreshScreen() {
	// clear screen
	editorScroll()
	ab := bytes.Buffer{}

	// hide cursor
	fmt.Fprint(&ab, "\x1b[?25l")

	// clear screen
	// fmt.Fprint(&ab, "\x1b[2J")

	// move cursor to top left
	fmt.Fprint(&ab, "\x1b[H")

	editorDrawRows(&ab)

	// reposition cursor
	//fmt.Fprint(&ab, "\x1b[H")
	fmt.Fprintf(&ab, "\x1b[%d;%dH", cfg.cy-cfg.rowOffset+1, cfg.cx+1)

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

func editorScroll() {
	if cfg.cy < cfg.rowOffset {
		cfg.rowOffset = cfg.cy
	}

	if cfg.cy >= cfg.rowOffset+cfg.screenRows {
		cfg.rowOffset = cfg.cy - cfg.screenRows + 1
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
			rowSize := len(cfg.rows[fileRow])
			if rowSize > cfg.screenCols {
				rowSize = cfg.screenCols
			}
			fmt.Fprint(ab, string(cfg.rows[fileRow][:rowSize]))
		}

		// clear to end of line
		fmt.Fprint(ab, "\x1b[K")

		if y < cfg.screenRows-1 {
			fmt.Fprint(ab, "\r\n")
		}
	}
}
