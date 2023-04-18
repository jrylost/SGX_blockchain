package server

import (
	"SGX_blockchain/src/accounts"
	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/db"
	"SGX_blockchain/src/utils"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const localhost string = `https://127.0.0.1`
const jsonContentType string = "application/json"

func HttpPost(url string, body string) (string, error) {
	resp, err := http.Post(url, jsonContentType, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	return string(result), err
}

func TestAccountsHandler(t *testing.T) {
	d := db.InitMemorydb()
	h := NewHandler(d)

	keypair, err := crypto.NewKeyPair()
	if err != nil {
		panic("Key generation wrong!")
	}

	account := &accounts.ExternalAccount{
		Id:      keypair.PubK,
		Count:   5,
		Nonce:   100,
		Balance: 2000,
	}

	accountInfo, _ := json.Marshal(account)
	d.TryLock(keypair.PubK)
	suc := d.Put(keypair.PubK, accountInfo)
	info, _ := d.Get(keypair.PubK)
	fmt.Println()
	fmt.Println(utils.EncodeBytesToHexStringWith0x(info), suc)
	fmt.Println(utils.EncodeBytesToHexStringWith0x(keypair.PubK))

	testserver := httptest.NewServer(http.HandlerFunc(h.AccountInfoHandler))
	accountsRequest := &AccountInfoRequest{
		Data: struct {
			Address string `json:"address"`
			Ts      int64  `json:"ts"`
		}{
			Address: utils.EncodeBytesToHexStringWith0x(keypair.PubK),
			Ts:      time.Now().UnixMilli(),
		},
		Signature: "",
	}
	defer testserver.Close()

	bodyBytes, err := utils.SignJsonWithData(accountsRequest, keypair)
	if err != nil {
		fmt.Println("Wrong accountsRequest")
	}

	buffer := bytes.NewBuffer(bodyBytes)

	resp, _ := http.Post(testserver.URL, jsonContentType, buffer)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func BenchmarkAccountInfoRequest(b *testing.B) {
	d := db.InitMemorydb()
	h := NewHandler(d)

	keypair, err := crypto.NewKeyPair()
	if err != nil {
		panic("Key generation wrong!")
	}

	account := &accounts.ExternalAccount{
		Id:      keypair.PubK,
		Count:   5,
		Nonce:   100,
		Balance: 2000,
	}

	accountInfo, _ := json.Marshal(account)
	d.TryLock(keypair.PubK)
	d.Put(keypair.PubK, accountInfo)
	//info, _ := d.Get(keypair.PubK)
	//fmt.Println()
	//fmt.Println(utils.EncodeBytesToHexStringWith0x(info), suc)
	//fmt.Println(utils.EncodeBytesToHexStringWith0x(keypair.PubK))

	testserver := httptest.NewServer(http.HandlerFunc(h.AccountInfoHandler))
	accountsRequest := &AccountInfoRequest{
		Data: struct {
			Address string `json:"address"`
			Ts      int64  `json:"ts"`
		}{
			Address: utils.EncodeBytesToHexStringWith0x(keypair.PubK),
			Ts:      time.Now().UnixMilli(),
		},
		Signature: "",
	}
	defer testserver.Close()

	bodyBytes, err := utils.SignJsonWithData(accountsRequest, keypair)
	if err != nil {
		fmt.Println("Wrong accountsRequest")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// do something
		utils.HttpPostWithClient(cli, testserver.URL, bodyBytes)

	}

}

var cli = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
