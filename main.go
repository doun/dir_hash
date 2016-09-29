package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"os"
)

func Hash_file(f string) (hash string, err error) {
	h := sha1.New()
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		return
	}
	io.Copy(h, file)
	finfo, err := os.Stat(f)
	if err != nil {
		return
	}
	return fmt.Sprintf("%s,%v,%x", finfo.Name(), finfo.Size(), h.Sum(nil)), nil
}

func main() {
	dir := flag.String("p", ".", "dir to hash")
	exclude := flag.String("x", "", "file names to exclude")
	saveto := flag.String("o", "hash", "file name to save hashed")
	fmt.Print(dir, exclude, saveto)
	fmt.Println("vim-go")
}
