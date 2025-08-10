package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := walkDir(out, path, printFiles, "")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}

func walkDir(out io.Writer, path string, printFiles bool, prefix string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.New("ошибка открытия директории")
	}
	items, err := file.ReadDir(0)
	if err != nil {
		return errors.New("ошибка чтения директории")
	}
	filtered := make([]os.DirEntry, 0, len(items))
	for _, item := range items {
		if item.Name() == ".DS_Store" {
			continue
		}
		if !printFiles && !item.IsDir() {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name() < filtered[j].Name()
	})

	for i, item := range filtered {
		isLast := i == len(filtered)-1
		var symbol string
		if item.Name() == "project" && printFiles {
			symbol += "\n"
		}
		if isLast {
			symbol += "└───"
		} else {
			symbol += "├───"
		}

		info, _ := item.Info()
		if item.IsDir() {
			fmt.Fprintf(out, "%s%s%s\n", prefix, symbol, item.Name())
			var newPrefix string
			if isLast {
				newPrefix = prefix + "\t"
			} else {
				newPrefix = prefix + "│\t"
			}
			walkDir(out, filepath.Join(path, item.Name()), printFiles, newPrefix)
		} else if printFiles {
			var nameAndInfo string
			if info.Size() == 0 {
				nameAndInfo = fmt.Sprintf("%s (empty)", item.Name())
			} else {
				nameAndInfo = fmt.Sprintf("%s (%db)", item.Name(), info.Size())
			}
			_, err := fmt.Fprintf(out, "%s%s%s\n", prefix, symbol, nameAndInfo)
			if err != nil {
				return errors.New("ошибка преобразования в строку")
			}
		}
	}
	return nil
}
