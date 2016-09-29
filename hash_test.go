package main

import (
	"fmt"
	"strings"
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

func Test_exclude_building(t *testing.T) {
	x := "*a, a*, *a*, cdd"
	contains := []string{"bcd", "bbd"}
	not_contains := []string{"aa", "abb", "cdd"}
	excludes := build_excludes(x)
	for _, p := range contains {
		if FileExcluded(p, excludes) {
			t.Fail()
		}
	}
	for _, p := range not_contains {
		if !FileExcluded(p, excludes) {
			t.Fail()
		}
	}
}

func Test_strings_prefix(t *testing.T) {
	fmt.Println(strings.HasPrefix("aa", "a"))
	if !strings.HasPrefix("a", "a") {
		t.Fail()
	}
}
