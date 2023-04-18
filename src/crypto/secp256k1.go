package crypto

import (
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

func EthereumSignatureToCompact(sig []byte) []byte {
	output := make([]byte, len(sig))
	output[0] = sig[len(sig)-1] + 27
	copy(output[1:], sig)
	return output
}

func CompactToEthereumSignature(sig []byte) []byte {
	output := make([]byte, len(sig))
	output[len(sig)-1] = sig[0] - 27
	copy(output[:len(sig)-1], sig[1:])
	return output
}

const (
	// asn1SequenceID is the ASN.1 identifier for a sequence and is used when
	// parsing and serializing signatures encoded with the Distinguished
	// Encoding Rules (DER) format per section 10 of [ISO/IEC 8825-1].
	asn1SequenceID = 0x30

	// asn1IntegerID is the ASN.1 identifier for an integer and is used when
	// parsing and serializing signatures encoded with the Distinguished
	// Encoding Rules (DER) format per section 10 of [ISO/IEC 8825-1].
	asn1IntegerID = 0x02
)

func EthereumSignatureToDER(sig []byte) []byte {
	canonR, canonS := make([]byte, 33), make([]byte, 33)
	copy(canonR[1:], sig[:32])
	copy(canonS[1:], sig[32:64])
	for len(canonR) > 1 && canonR[0] == 0x00 && canonR[1]&0x80 == 0 {
		canonR = canonR[1:]
	}
	for len(canonS) > 1 && canonS[0] == 0x00 && canonS[1]&0x80 == 0 {
		canonS = canonS[1:]
	}

	// Total length of returned signature is 1 byte for each magic and length
	// (6 total), plus lengths of R and S.
	totalLen := 6 + len(canonR) + len(canonS)
	b := make([]byte, 0, totalLen)
	b = append(b, asn1SequenceID)
	b = append(b, byte(totalLen-2))
	b = append(b, asn1IntegerID)
	b = append(b, byte(len(canonR)))
	b = append(b, canonR...)
	b = append(b, asn1IntegerID)
	b = append(b, byte(len(canonS)))
	b = append(b, canonS...)
	return b
}

func VerifyHashSignature(sig, hash, pubK []byte) bool {
	pubKey, _ := btcec.ParsePubKey(pubK)
	signature, err := ecdsa.ParseDERSignature(sig)
	if err != nil {
		fmt.Println(err)
	}
	return signature.Verify(hash, pubKey)
}

func VerifyMessageSignature(sig, pubK []byte, message ...[]byte) bool {
	hash := Keccak256(message...)
	return VerifyHashSignature(sig, hash, pubK)
}

type KeyPair struct {
	PriK       []byte
	PubK       []byte
	PrivateKey *btcec.PrivateKey
	PublicKey  *btcec.PublicKey
}

func Initialize(pk []byte) *KeyPair {
	k := &KeyPair{PriK: pk}
	k.PrivateKey, k.PublicKey = btcec.PrivKeyFromBytes(pk)
	k.PubK = k.PublicKey.SerializeCompressed()
	return k
}

func NewKeyPair() (*KeyPair, error) {
	k := &KeyPair{}
	var err error
	k.PrivateKey, err = btcec.NewPrivateKey()
	k.PublicKey = k.PrivateKey.PubKey()
	k.PubK = k.PublicKey.SerializeCompressed()
	k.PriK = k.PrivateKey.Serialize()
	return k, err
}

func (k *KeyPair) SignEthereumHash(hash []byte) []byte {
	sig, err := ecdsa.SignCompact(k.PrivateKey, hash, false)
	if err != nil {
		fmt.Println("Sign error")
	}
	//sig := ecdsa.Sign(k.PrivateKey, hash)
	return append(sig, sig[0]-0x1c)[1:]
}

func (k *KeyPair) SignMessage(message ...[]byte) []byte {
	hash := Keccak256(message...)
	return k.SignEthereumHash(hash)
}

func ToCompressedPubKey(pubkey []byte) []byte {
	if len(pubkey) == 33 {
		return pubkey
	} else {
		p, _ := btcec.ParsePubKey(pubkey)
		return p.SerializeCompressed()
	}
}

//
//func test() {
//	newAEAD := func(key []byte) (cipher.AEAD, error) {
//		block, err := aes.NewCipher(key)
//		if err != nil {
//			return nil, err
//		}
//		return cipher.NewGCM(block)
//	}
//
//	// Decode the hex-encoded pubkey of the recipient.
//	pubKeyBytes, err := hex.DecodeString("04115c42e757b2efb7671c578530ec191a1" +
//		"359381e6a71127a9d37c486fd30dae57e76dc58f693bd7e7010358ce6b165e483a29" +
//		"21010db67ac11b1b51b651953d2") // uncompressed pubkey
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
//
//	// Using ECDHE, derive a shared symmetric key for encryption of the plaintext.
//	cipherKey := sha256.Sum256(secp256k1.GenerateSharedSecret(ephemeralPrivKey, pubKey))
//
//	// Seal the message using an AEAD.  Here we use AES-256-GCM.
//	// The ephemeral public key must be included in this message, and becomes
//	// the authenticated data for the AEAD.
//	//
//	// Note that unless a unique nonce can be guaranteed, the ephemeral
//	// and/or shared keys must not be reused to encrypt different messages.
//	// Doing so destroys the security of the scheme.  Random nonces may be
//	// used if XChaCha20-Poly1305 is used instead, but the message must then
//	// also encode the nonce (which we don't do here).
//	plaintext := []byte("test message")
//	aead, err := newAEAD(cipherKey[:])
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	nonce := make([]byte, aead.NonceSize())
//	ciphertext := make([]byte, 4+len(ephemeralPubKey))
//	binary.LittleEndian.PutUint32(ciphertext, uint32(len(ephemeralPubKey)))
//	copy(ciphertext[4:], ephemeralPubKey)
//	ciphertext = aead.Seal(ciphertext, nonce, plaintext, ephemeralPubKey)
//
//	// The remainder of this example is performed by the recipient on the
//	// ciphertext shared by the sender.
//
//	// Decode the hex-encoded private key.
//	pkBytes, err := hex.DecodeString("a11b0a4e1a132305652ee7a8eb7848f6ad" +
//		"5ea381e3ce20a2c086a2e388230811")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	privKey := secp256k1.PrivKeyFromBytes(pkBytes)
//
//	// Read the sender's ephemeral public key from the start of the message.
//	// Error handling for inappropriate pubkey lengths is elided here for
//	// brevity.
//	pubKeyLen := binary.LittleEndian.Uint32(ciphertext[:4])
//	senderPubKeyBytes := ciphertext[4 : 4+pubKeyLen]
//	senderPubKey, err := secp256k1.ParsePubKey(senderPubKeyBytes)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	// Derive the key used to seal the message, this time from the
//	// recipient's private key and the sender's public key.
//	recoveredCipherKey := sha256.Sum256(secp256k1.GenerateSharedSecret(privKey, senderPubKey))
//
//	// Open the sealed message.
//	aead, err = newAEAD(recoveredCipherKey[:])
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	nonce = make([]byte, aead.NonceSize())
//	recoveredPlaintext, err := aead.Open(nil, nonce, ciphertext[4+pubKeyLen:], senderPubKeyBytes)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(string(recoveredPlaintext))
//}
