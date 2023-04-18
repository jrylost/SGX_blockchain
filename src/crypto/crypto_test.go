package crypto

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAES(t *testing.T) {
	pk := generateAESKey()
	k := initialize(pk)
	plain := make([]byte, 4, 16)
	iv := make([]byte, 16)
	copy(plain, []byte("abcd"))
	rand.Read(iv)
	cipherText, _ := k.Encrypt(plain, iv)
	fmt.Println(hex.EncodeToString(cipherText))
	plainText, _ := k.Decrypt(cipherText)
	if !bytes.Equal(plainText, []byte("abcd")) {
		t.Fatalf("AES wrong!")
	}
	fmt.Println(string(plainText))
}

const benchmarkTestSliceLength = 1024 * 1024 * 10

func Benchmark(b *testing.B) {
	b.Run("BenchmarkAESwithPadding", BenchmarkAESwithPadding)
	b.Run("BenchmarkAESwithoutPadding", BenchmarkAESwithoutPadding)

}

func BenchmarkAESwithPadding(b *testing.B) {
	pk := generateAESKey()
	k := initialize(pk)
	plain := bytes.Repeat([]byte("a"), benchmarkTestSliceLength)
	plain = plain[:benchmarkTestSliceLength-AESBlockSize]
	//copy(plain, []byte("abcd"))
	//fmt.Println(plain[:20])
	iv := generateIV()
	b.ResetTimer()
	cipherText, _ := k.Encrypt(plain, iv)
	//fmt.Println(cipherText)
	plainText, _ := k.Decrypt(cipherText)
	b.StopTimer()
	//fmt.Println(len(plainText))
	if !bytes.Equal(plainText, plain) {
		b.Fatalf("AES wrong!")
	}
}

func BenchmarkAESwithoutPadding(b *testing.B) {
	//b.ReportAllocs()
	pk := generateAESKey()
	k := initialize(pk)
	plain := bytes.Repeat([]byte("a"), benchmarkTestSliceLength+1)
	//copy(plain, []byte("abcd"))
	iv := generateIV()
	b.ResetTimer()
	cipherText, _ := k.Encrypt(plain, iv)
	//b.ReportAllocs()
	plainText, _ := k.Decrypt(cipherText)
	b.StopTimer()
	if !bytes.Equal(plainText, plain) {
		b.Fatalf("AES wrong!")
	}
}

func BenchmarkAES(b *testing.B) {
	pk := generateAESKey()
	k := initialize(pk)
	plain := bytes.Repeat([]byte("a"), benchmarkTestSliceLength)
	plain = plain[:benchmarkTestSliceLength-AESBlockSize]
	//plain := bytes.Repeat([]byte("a"), benchmarkTestSliceLength+1)
	//plain = plain[:benchmarkTestSliceLength-AESBlockSize]

	iv := generateIV()
	for i := 0; i < b.N; i++ {
		cipherText, _ := k.Encrypt(plain, iv)
		cipherText[0] = 0x22
		//b.ReportAllocs()
		//k.Decrypt(cipherText)
	}
}

func BenchmarkAESParallel(b *testing.B) {
	pk := generateAESKey()
	k := initialize(pk)
	plain := bytes.Repeat([]byte("a"), benchmarkTestSliceLength)
	plain = plain[:benchmarkTestSliceLength-AESBlockSize]
	//plain := bytes.Repeat([]byte("a"), benchmarkTestSliceLength+1)
	//plain = plain[:benchmarkTestSliceLength-AESBlockSize]

	iv := generateIV()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// do something
			//fmt.Sprint("代码轶事")
			cipherText, _ := k.Encrypt(plain, iv)
			cipherText[0] = 0x22
		}
	})
	//for i := 0; i < b.N; i++ {
	//	//b.ReportAllocs()
	//	//k.Decrypt(cipherText)
	//}
}

//func TestSecp256k1Package(t *testing.T) {
//	pubKeyBytes, err := hex.DecodeString("04115c42e757b2efb7671c578530ec191a1" +
//		"359381e6a71127a9d37c486fd30dae57e76dc58f693bd7e7010358ce6b165e483a29" +
//		"21010db67ac11b1b51b651953d2")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	pubKey, err := secp256k1.ParsePubKey(pubKeyBytes)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// Derive an ephemeral public/private keypair for performing ECDHE with
//	// the recipient.
//	ephemeralPrivKey, err := secp256k1.GeneratePrivateKey()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	ephemeralPubKey := ephemeralPrivKey.PubKey().SerializeCompressed()
//}

func TestName(t *testing.T) {
	a := make([]byte, 5, 8)
	copy(a, []byte("abcde"))
	c := a[3:]
	fmt.Println(c)
	fmt.Println(cap(c), len(c))
}
