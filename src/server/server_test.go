package server

import (
	"SGX_blockchain/src/accounts"
	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/db"
	"SGX_blockchain/src/utils"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
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
	result, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	return string(result), err
}

func TestAccountsInfoHandler(t *testing.T) {
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
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func TestFileStoreHandler(t *testing.T) {
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

	testserver := httptest.NewServer(http.HandlerFunc(h.FileStoreHandler))

	filecontent := []byte("this is file content!")
	filehash := crypto.Keccak256(filecontent)
	filecontentbase64 := base64.StdEncoding.EncodeToString(filecontent)

	filestorerequest := &FileStoreRequest{
		Data: struct {
			Content  string `json:"content"`
			FileHash string `json:"fileHash"`
			From     string `json:"from"`
			Nonce    int64  `json:"nonce"`
			Ts       int64  `json:"ts"`
		}{
			Content:  filecontentbase64,
			FileHash: utils.EncodeBytesToHexStringWith0x(filehash),
			From:     utils.EncodeBytesToHexStringWith0x(keypair.PubK),
			Nonce:    1,
			Ts:       time.Now().UnixMilli(),
		},
		Signature: "",
	}
	fmt.Println(filecontentbase64, "base64")
	defer testserver.Close()

	bodyBytes, err := utils.SignJsonWithData(filestorerequest, keypair)
	if err != nil {
		fmt.Println("Wrong accountsRequest")
	}

	buffer := bytes.NewBuffer(bodyBytes)

	resp, _ := http.Post(testserver.URL, jsonContentType, buffer)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	//fmt.Println(d)

	testserver2 := httptest.NewServer(http.HandlerFunc(h.FileRetrieveHandler))
	fileretrieverequest := &FileRetrieveRequest{
		Data: struct {
			From     string `json:"from"`
			FileHash string `json:"fileHash"`
			Ts       int64  `json:"ts"`
		}{
			From:     utils.EncodeBytesToHexStringWith0x(keypair.PubK),
			FileHash: utils.EncodeBytesToHexStringWith0x(filehash),
			Ts:       time.Now().UnixMilli(),
		},
		Signature: "",
	}
	bodyBytes2, err := utils.SignJsonWithData(fileretrieverequest, keypair)
	if err != nil {
		fmt.Println("Wrong accountsRequest")
	}

	buffer2 := bytes.NewBuffer(bodyBytes2)
	resp2, _ := http.Post(testserver2.URL, jsonContentType, buffer2)
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	fmt.Println(string(body2))
	base64str := gjson.GetBytes(body2, "transaction.content").String()
	fcontent, _ := base64.StdEncoding.DecodeString(base64str)

	fmt.Println(string(fcontent))
	defer testserver2.Close()

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
