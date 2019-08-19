package main

import (
	"bufio"
	"fmt"
	"os"
)

func editorOpen(fileName string) error {

	r, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	// TK : WOrking
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	defText := "Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World Hello World"
	row := erow{}
	for _, runeVal := range defText {
		row = append(row, runeVal)
	}

	cfg.rows = []erow{row}
	return nil
}
