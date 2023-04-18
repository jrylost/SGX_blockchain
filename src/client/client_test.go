package client

import (
	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/utils"
	"crypto/tls"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"testing"
	"time"
)

var pair, _ = crypto.NewKeyPair()
var c = NewClient(&tls.Config{InsecureSkipVerify: true}, pair)
var cli = &http.Client{Transport: &http.Transport{TLSClientConfig: c.tls}}

func TestRequestForAccountInfo(t *testing.T) {
	c.RequestForAccountInfo(pair)
}

func BenchmarkClient_RequestForAccountInfo(b *testing.B) {
	//tlsConfig := &tls.Config{InsecureSkipVerify: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.HttpPostWithClient(cli, ServerURL+"/report", nil)
	}

	//b.RunParallel(func(pb *testing.PB) {
	//	for pb.Next() {
	//		go utils.HttpPost(tlsConfig, ServerURL+"/", nil)
	//	}
	//})
}

func TestBasic(t *testing.T) {
	//tlsConfig := &tls.Config{InsecureSkipVerify: true}
	//resp := utils.HttpPost(tlsConfig, ServerURL+"/report", nil)
	////if err != nil {
	////	fmt.Println(err)
	////} else {
	//var b int64
	//fmt.Scan("%d", &b)
	//fmt.Println(utils.EncodeBytesToHexStringWith0x(resp))
	////}

	for _, i := range Response {
		time.Sleep(317 * time.Millisecond)
		if gjson.Valid(i) {
			i, _ = sjson.Set(i, "data.ts", time.Now().Unix())

		}
		fmt.Println(i)
	}
}
