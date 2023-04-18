package utils

import (
	"crypto/rsa"
	"github.com/edgelesssys/ego/enclave"
	"time"
)

func CreateCertificate() ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "SGX_Blockchain"},
		NotAfter:     time.Now().AddDate(10, 0, 0),
		DNSNames:     []string{"localhost"},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}

func GetRemoteReport(hash []byte) []byte {
	report, err := enclave.GetRemoteReport(hash[:])
	if err != nil {
		fmt.Println(err)
	}
	return report
}
