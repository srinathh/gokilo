package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func editorOpen(fileName string) error {

	r, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	cfg.rows = []erow{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		editorAppendRow(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	cfg.fileName = fileName
	cfg.dirty = false
	return nil
}

func editorRowsToString() string {
	var sb strings.Builder

	for _, rows := range cfg.rows {
		sb.WriteString(string(rows.chars))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func editorSave() {
	if cfg.fileName == "" {
		return
	}

	fil, err := os.Create(cfg.fileName)
	if err != nil {
		editorSetStatusMsg("ERROR creating file: %s: %s", err, cfg.fileName)
		return
	}
	defer fil.Close()

	if _, err := fmt.Fprint(fil, editorRowsToString()); err != nil {
		editorSetStatusMsg("ERROR writing to file: %s: %s", err, cfg.fileName)
		return
	}

	if err = fil.Close(); err != nil {
		editorSetStatusMsg("ERROR closing written file: %s: %s", err, cfg.fileName)
		return
	}

	editorSetStatusMsg("SAVED to file: %s", cfg.fileName)
	cfg.dirty = false

}
