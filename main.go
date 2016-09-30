package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	BlockSize = 64 * 64
)

func Hash_file(path string) (hash string, err error) {
	h := sha1.New()
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return
	}
	buf := make([]byte, BlockSize)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			h.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

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
	dir := flag.String("d", ".", "dir to be hashed")
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

	rst_channel := make(chan string, 10)
	n_found := 0
	close_channel := make(chan bool, 1)

	err = filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
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
			n_found++
			go func() {
				h, err := Hash_file(path)
				if err == nil {
					rst_channel <- h + "\r\n"
				}
			}()
		}
		return nil
	})

	if err != nil {
		flag.Usage()
	}

	go func() {
		for {
			rst := <-rst_channel
			h_file.WriteString(rst)
			n_found--
			if n_found == 0 {
				close_channel <- true
				break
			}
		}
	}()

	signal := <-close_channel
	if signal {
		close(rst_channel)
		return
	}
}
