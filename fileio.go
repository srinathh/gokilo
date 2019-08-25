package main

import (
	"bufio"
	"os"
)

func editorOpen(fileName string) error {

	r, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	cfg.rows = []erow{}

	// TK : Working
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		row := erow{
			chars: []rune(scanner.Text()),
		}
		cfg.rows = append(cfg.rows, row)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
