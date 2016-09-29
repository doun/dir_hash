package main

import (
	"testing"
)

func Test_hash_file(t *testing.T) {
	s, err := Hash_file("./testdata/file1")
	if err != nil {
		t.Fatal(err)
	}
	if "file1,9,f64471bb8418b892618a5f3a73bb59ebadb9181c" != s {
		t.Fatal("got:", s, "expecting", "file1,9,f64471bb8418b892618a5f3a73bb59ebadb9181c")
	}
}
