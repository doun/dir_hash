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

func FileExcluded(path string, x []exclude) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	name := parts[len(parts)-1]
	for _, s := range x {
		if s.Exclude(name) {
			return true
		}
	}
	return false
}

func build_excludes(x string) []exclude {
	var excludes []exclude
	xarray := strings.Split(x, ",")
	for _, s := range xarray {
		s = strings.Trim(s, " ")
		if len(s) == 0 {
			continue
		}
		exp := strings.Trim(s, "*")
		if strings.HasPrefix(s, "*") {
			excludes = append(excludes, exclude{exp, func(p, e string) bool {
				return strings.HasPrefix(p, e)
			}})
		}
		if strings.HasSuffix(s, "*") {
			excludes = append(excludes, exclude{exp, func(p, e string) bool {
				return strings.HasSuffix(p, e)
			}})
		}
		if strings.Index(s, "*") < 0 {
			excludes = append(excludes, exclude{exp, func(p, e string) bool {
				return p == e
			}})
		}
	}
	return excludes
}

type exclude struct {
	exp   string
	check func(string, string) bool
}

func (self *exclude) Exclude(path string) bool {
	return self.check(path, self.exp)
}

func main() {
	dir := flag.String("p", ".", "dir to hash")
	x := flag.String("x", "", "file/dir name to exclude, seperate with ','")
	saveto := flag.String("o", "", "file name to save hashed info, will be truncked if exist!!")
	flag.Parse()

	excludes := build_excludes(*x)

	h_file, err := os.Create(*saveto)
	defer h_file.Close()
	if err != nil {
		flag.Usage()
		return
	}
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
				h_file.WriteString(h + "\r\n")
			}
		}
		return nil
	})
}
