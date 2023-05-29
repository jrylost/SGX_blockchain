package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

func TestKeccak256Hash3(t *testing.T) {
	//sss, _ := hex.DecodeString("0450863ad64a87ae8a2fe83c1af1a8403cb53f53e486d8511dad8a04887e5b23522cd470243453a299fa9e77237716103abc11a1df38855ed6f2ee187e9c582ba6")
	sss, _ := hex.DecodeString("dfa13518ff965498743f3a01439dd86bc34ff9969c7a3f0430bbf8865734252953c9884af787b2cadd45f92dff2b81e21cfdf98873e492e5fdc07e9eb67ca74d")
	fmt.Println(hex.EncodeToString(Keccak256(sss)))
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

func TestSHA256Hash(t *testing.T) {
	msg := []byte("")
	exp, _ := hex.DecodeString("a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := sha256(in); return h[:] }, msg, exp)
}
