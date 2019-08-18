package main

import (
	"bytes"
	"fmt"
	"os"
)

func editorRefreshScreen() {
	// clear screen
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
	fmt.Fprintf(&ab, "\x1b[%d;%dH", cfg.cy+1, cfg.cx+1)

	// show cursor
	fmt.Fprint(&ab, "\x1b[?25h")

	os.Stdout.Write(ab.Bytes())

}

func editorDrawRows(ab *bytes.Buffer) {
	for y := 0; y < cfg.screenRows; y++ {

		if y == cfg.screenRows/3 {
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

		// clear to end of line
		fmt.Fprint(ab, "\x1b[K")

		if y < cfg.screenRows-1 {
			fmt.Fprint(ab, "\r\n")
		}
	}
}
