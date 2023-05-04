package db

import (
	"SGX_blockchain/src/vm"
	"bytes"
	"fmt"
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

func TestMemorydbcontext(t *testing.T) {
	d := InitMemorydb()

	si := vm.StorageInterface{
		StringStorage: map[string]string{
			"aaaccc": "bbbccc",
		},
	}
	cdb := d.contextdb
	cdb["testaaa"] = si
	//d.contextdb = contextdb

	fmt.Println(d)
}
