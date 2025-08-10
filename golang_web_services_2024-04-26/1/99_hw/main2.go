package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	out := os.Stdout
	path := os.Args[0]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := walkDir(out, path, printFiles)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}

func walkDir(out io.Writer, path string, printFiles bool) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.New("ошибка открытия дериктории")
	}
	names, _ := file.Readdirnames(0)
	for _, name := range names {
		fmt.Println("Имена", name)
	}
	return nil
}
