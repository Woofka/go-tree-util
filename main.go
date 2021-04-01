package main

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	lineSpace  = "    "
	lineTab    = "│   "
	lineMiddle = "├───"
	lineLast   = "└───"
)

func printTree(out io.Writer, path string, printFiles bool, prefix string) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	if !printFiles {
		n := 0
		for _, x := range files {
			if x.IsDir() {
				files[n] = x
				n++
			}
		}
		files = files[:n]
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	for i, fileInfo := range files {
		postfix := ""
		if !fileInfo.IsDir() {
			if fileSize := fileInfo.Size(); fileSize == 0 {
				postfix = " (empty)"
			} else {
				postfix = " (" + strconv.Itoa(int(fileSize)) + "b)"
			}

		}
		var line string
		isLastLine := i == len(files)-1
		if isLastLine {
			line = lineLast
		} else {
			line = lineMiddle
		}
		_, err := out.Write([]byte(prefix + line + fileInfo.Name() + postfix + "\n"))
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			nextPrefix := prefix
			if isLastLine {
				nextPrefix += lineSpace
			} else {
				nextPrefix += lineTab
			}
			err := printTree(out, filepath.Join(path, fileInfo.Name()), printFiles, nextPrefix)
			if err != nil {
				return err
			}
		}
	}
	return dir.Close()
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := printTree(out, path, printFiles, "")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("Usage: go run main.go <PATH> [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
