package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func checkhash(t *testing.T, name string, f func([]byte) []byte, msg, exp []byte) {
	sum := f(msg)
	if !bytes.Equal(exp, sum) {
		t.Fatalf("hash %s mismatch: want: %x have: %x", name, exp, sum)
	}
}

// Sanity checks.
func TestKeccak256Hash(t *testing.T) {
	msg := []byte("abc")
	exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := Keccak256(in); return h[:] }, msg, exp)
}

// Sanity checks.
func TestKeccak256Hash2(t *testing.T) {
	msg := []byte("a")
	msg2 := []byte("b")
	msg3 := []byte("c")
	exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := Keccak256(msg, msg2, msg3); return h[:] }, msg, exp)
}

// BenchmarkKeccakHash
// @Description 测试keccak哈希函数性能
// @Author jerry 2022-09-24 22:24:56
// @Param b
func BenchmarkKeccakHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Keccak256([]byte("a"))
	}
}
