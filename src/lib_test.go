package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGrep(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "sec")
	if err != nil {
		t.Error(err)
	}
	f, err := os.Create(dir + "/sec_unit_test.txt")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	_, err = f.WriteString("public_key \n private_key")

	if err != nil {
		t.Error(err)
	}

	findings, err := FindWords(dir)

	if len(findings) != 2 {
		t.Error("Grep is failed in finding words")
	}
	f.Close()
	os.RemoveAll(dir)
}
