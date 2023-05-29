package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomBytes(length int64) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func GenerateRandomHexString(length int64) string {
	b := GenerateRandomBytes(length)
	s := hex.EncodeToString(b)
	return s
}

func GenerateRandomHexStringWith0x(length int64) string {
	b := GenerateRandomBytes(length)
	s := hex.EncodeToString(b)
	return "0x" + s
}

func EncodeBytesToHexStringWith0x(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func DecodeHexStringToBytesWith0x(s string) []byte {
	val, _ := hex.DecodeString(s[2:])
	return val
}
