package main

import (
	"fmt"
	"gokilo/rawmode"
	"gokilo/terminal"
	"os"
	"unicode"
)

func main() {

	bkup, err := rawmode.Enable()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	defer rawmode.Restore(bkup)

	fmt.Fprint(os.Stdout, "Press q to quit\r\n")
	printKey()
}

func printKey() {

forLoop:
	for {
		k, err := terminal.ReadKey()

		switch {

		case err == terminal.ErrNoInput:
			continue forLoop

		case err != nil:
			fmt.Fprintf(os.Stderr, "%s\r\n", err)
			break forLoop

		case k.Special != terminal.KeyNoSpl:
			fmt.Printf("Special key pressed: %d\r\n", k.Special)
			continue forLoop

		case unicode.ToLower(k.Regular) == 'q':
			break forLoop

		case k.Regular >= 32 && k.Regular != 127:
			fmt.Fprintf(os.Stdout, "Got Key: %d %c\r\n", k.Regular, k.Regular)
			continue forLoop

		default:
			fmt.Fprintf(os.Stdout, "Got Key: %d\r\n", k.Regular)

		}
	}
}
