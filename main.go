package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Hash_file(path string) (hash string, err error) {
	h := sha1.New()
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}
	io.Copy(h, file)
	finfo, err := os.Stat(path)
	if err != nil {
		return
	}
	return fmt.Sprintf("%s,%v,%x", finfo.Name(), finfo.Size(), h.Sum(nil)), nil
}

func FileExcluded(path string, x []string) bool {
	return false
}

func main() {
	dir := flag.String("p", ".", "dir to hash")
	exclude := flag.String("x", "", "file/dir name to exclude, seperate with ','")
	saveto := flag.String("o", "hash", "file name to save hashed info, will be truncked if exist!!")
	flag.Parse()

	h_file, err := os.Create(*saveto)
	if err != nil {
		return
	}
	defer h_file.Close()

	excludes := strings.Split(*exclude, ",")

	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if FileExcluded(path, excludes) {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if !info.IsDir() {
			h, err := Hash_file(path)
			if err == nil {
				h_file.WriteString(h)
			}
		}
		return nil
	})
	fmt.Print(*dir, exclude, saveto)
}
