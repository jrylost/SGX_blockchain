package crypto

import (
	"encoding/hex"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"testing"
)

var (
	testprivk, _   = hex.DecodeString("d9b29c4dc2b3202d8b43ae6677ac02e754a1217f94a193c53f21e89d26d06685")
	testmsg, _     = hex.DecodeString("ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testsig, _     = hex.DecodeString("90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey, _  = hex.DecodeString("04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc, _ = hex.DecodeString("02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
)

func TestEcrecover(t *testing.T) {
	testSig2 := EthereumSignatureToCompact(testsig)
	//fmt.Println(hex.EncodeToString(testSig2))
	pubKey, _, err := ecdsa.RecoverCompact(testSig2, testmsg)
	if err != nil {
		t.Fatalf("recover error: %s", err)
	}
	uncompressedPubK := hex.EncodeToString(pubKey.SerializeCompressed())
	fmt.Println(uncompressedPubK)
}

// 测试压缩公钥和非压缩公钥的签名验证
func TestVerifySig(t *testing.T) {
	//fmt.Println(hex.EncodeToString(testsig))
	sig := EthereumSignatureToDER(testsig)
	//fmt.Println(hex.EncodeToString(testsig))
	bool1 := VerifyHashSignature(sig, testmsg, testpubkeyc)
	bool2 := VerifyHashSignature(sig, testmsg, testpubkey)
	//fmt.Println(hex.EncodeToString(sig))
	//fmt.Println(hex.EncodeToString(testmsg))
	//fmt.Println(hex.EncodeToString(testpubkey))
	if !(bool1 && bool2) {
		t.Fatalf("verification wrong!")
	}
}

func TestSign(t *testing.T) {
	k, err := NewKeyPair()
	if err != nil {
		t.Fatalf("generate key wrong!")
	}

	fmt.Println(EncodeBytesToHexStringWith0x(k.PrivateKey.Serialize()))
	fmt.Println(EncodeBytesToHexStringWith0x(k.PublicKey.SerializeCompressed()))
	fmt.Println(EncodeBytesToHexStringWith0x(k.PublicKey.SerializeUncompressed()))
	msg := [][]byte{[]byte("abc"), []byte("bcd")}
	//msg := "ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008"
	sig := k.SignMessage(msg...)
	fmt.Println(hex.EncodeToString(sig))
	b := VerifyMessageSignature(EthereumSignatureToDER(sig), k.PubK, msg...)
	if !b {
		t.Fatalf("wrong signature!")
	}
}

func TestSign2(t *testing.T) {
	k, _ := Initialize(testprivk)
	//if err != nil {
	//	t.Fatalf("generate key wrong!")
	//}

	fmt.Println(hex.EncodeToString(k.PriK))
	fmt.Println(hex.EncodeToString(k.PublicKey.SerializeUncompressed()))
	fmt.Println(hex.EncodeToString(k.PublicKey.SerializeCompressed()))
	fmt.Println(hex.EncodeToString(Keccak256(k.PublicKey.SerializeUncompressed()[1:])))
	fmt.Println(len(k.PublicKey.SerializeUncompressed()[1:]))
	fmt.Println(hex.EncodeToString(k.Address))
	//for i := 0; i < 100; i++ {
	//
	//fmt.Println(hex.EncodeToString(Keccak256(k.PublicKey.SerializeUncompressed()[i:])))
	//}

	fmt.Println("私钥：", EncodeBytesToHexStringWith0x(k.PrivateKey.Serialize()))
	fmt.Println("压缩公钥：", EncodeBytesToHexStringWith0x(k.PublicKey.SerializeCompressed()))
	fmt.Println("公钥：", EncodeBytesToHexStringWith0x(k.PublicKey.SerializeUncompressed()))
	msg := [][]byte{[]byte("abcbcd")}
	//msg := "ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008"
	sig := k.SignMessage(msg...)
	fmt.Println("签名：", EncodeBytesToHexStringWith0x(sig))
	b := VerifyMessageSignature(EthereumSignatureToDER(sig), k.PubK, msg...)
	if !b {
		t.Fatalf("wrong signature!")
	}
}

func TestSign3(t *testing.T) {
	k, err := NewKeyPair()
	if err != nil {
		t.Fatalf("generate key wrong!")
	}

	fmt.Println(EncodeBytesToHexStringWith0x(k.PrivateKey.Serialize()))
	fmt.Println(EncodeBytesToHexStringWith0x(k.PublicKey.SerializeCompressed()))
	fmt.Println(EncodeBytesToHexStringWith0x(k.PublicKey.SerializeUncompressed()))
	msg := [][]byte{[]byte("abc"), []byte("bcd")}
	//msg := "ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008"
	sig := k.SignMessage(msg...)
	fmt.Println(hex.EncodeToString(sig))
	b := VerifyMessageSignature(EthereumSignatureToDER(sig), k.PubK, msg...)
	if !b {
		t.Fatalf("wrong signature!")
	}
}

func BenchmarkVerifyHashSignature(b *testing.B) {
	sig := EthereumSignatureToDER(testsig)
	for i := 0; i < b.N; i++ {
		bool1 := VerifyHashSignature(sig, testmsg, testpubkeyc)
		if !bool1 {
			b.Fatalf("Wrong")
		}
	}
}
