package utils

import (
	"SGX_blockchain/src/crypto"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

func SignJsonWithData(structBody interface{}, k *crypto.KeyPair) ([]byte, error) {
	jsonBytes, err := json.Marshal(structBody)
	if err != nil {
		return []byte(""), err
	}
	data := gjson.GetBytes(jsonBytes, "data")
	if !data.Exists() {
		return []byte(""), errors.New("No data field!")
	}
	body := []byte(data.String())
	sig := k.SignMessage(body)

	bodyBytes, err := sjson.SetBytes(jsonBytes, "signature", EncodeBytesToHexStringWith0x(sig))
	if err != nil {
		return []byte(""), errors.New("Set signature error!")
	}
	return bodyBytes, nil

}
