package db

import (
	"bytes"
	"testing"
)

func TestMemorydb(t *testing.T) {
	d := InitMemorydb()
	d.Put([]byte("abcde"), []byte("aaaaa"))
	if value, b := d.Get([]byte("abcde")); !(b || !bytes.Equal(value, []byte("aaaaa"))) {
		t.Fatalf("Error")
	}
	if _, b := d.Get([]byte("jjkjk")); b {
		t.Fatalf("Error")
	}
}
