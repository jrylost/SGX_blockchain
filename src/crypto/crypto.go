package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const AESBlockSize = 16
const AESKeySize = 32

func generateIV() []byte {
	iv := make([]byte, AESBlockSize)
	rand.Read(iv)
	return iv
}

func generateAESKey() []byte {
	key := make([]byte, AESKeySize)
	rand.Read(key)
	return key
}

func (k *AESKey) EncryptWithReader(reader io.Reader, length int) []byte {
	lengthAfterPadding := length + AESBlockSize - (length % AESBlockSize)
	encrypted := make([]byte, AESBlockSize+length, AESBlockSize+lengthAfterPadding)
	iv := generateIV()
	copy(encrypted, iv)
	encryptedBytes := encrypted[AESBlockSize:]
	reader.Read(encryptedBytes)
	k.Encrypt(encryptedBytes, iv)
	return encrypted
}

func (k *AESKey) DecryptWithReader(reader io.Reader) []byte {
	decryptedBytes, _ := io.ReadAll(reader)
	decrypted, _ := k.Decrypt(decryptedBytes)
	return decrypted
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	//if cap(ciphertext)%blockSize != 0 {
	//	//panic("padding error!")
	//	log.Println("Padding length not reserved!")
	//}
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AESKey AES key
type AESKey struct {
	pk     []byte
	cipher cipher.Block
}

func initialize(pk []byte) *AESKey {
	key := &AESKey{
		pk: pk,
	}
	if len(key.pk) != AESKeySize {
		panic("private key length is not 32!")
	}
	key.cipher, _ = aes.NewCipher(pk)
	return key
}

// Encrypt :AES Encryption
func (k *AESKey) Encrypt(origData, iv []byte) ([]byte, error) {
	//add := &origData[0]
	origData = PKCS7Padding(origData, AESBlockSize)
	//if &origData[0] != add {
	//	fmt.Println("Padding Address changed!")
	//	fmt.Println(add, &origData[0])
	//}
	blockMode := cipher.NewCBCEncrypter(k.cipher, iv)
	encrypted := make([]byte, len(origData)+AESBlockSize)
	copy(encrypted, iv)
	encryptedPart := encrypted[AESBlockSize:]
	blockMode.CryptBlocks(encryptedPart, origData)
	return encrypted, nil
}

// Decrypt :AES Decryption
func (k *AESKey) Decrypt(encrypted []byte) ([]byte, error) {
	iv := encrypted[:AESBlockSize]
	blockMode := cipher.NewCBCDecrypter(k.cipher, iv)
	origData := make([]byte, len(encrypted)-AESBlockSize)
	blockMode.CryptBlocks(origData, encrypted[AESBlockSize:])
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
