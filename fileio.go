package main

func editorOpen() error {

	defText := "Hello World"
	row := erow{}
	for _, runeVal := range defText {
		row = append(row, runeVal)
	}

	cfg.rows = []erow{row}
	return nil
}
