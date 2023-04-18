package utils

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

const jsonContentType string = "application/json"

// HttpPost
// @Description  不可复用的HTTP Post。
// @Author jerry 2022-09-24 06:43:59
// @Param tlsConfig
// @Param url
// @Param reqBody
// @Return []byte
func HttpPost(tlsConfig *tls.Config, url string, reqBody []byte) []byte {
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	buffer := bytes.NewBuffer(reqBody)
	resp, err := client.Post(url, jsonContentType, buffer)
	//certInfo := resp.TLS.PeerCertificates[0]
	//fmt.Println("过期时间:", certInfo.NotAfter)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return respBody
}

// HttpPostWithClient
// @Description 可复用的HTTP Post，有客户端。
// @Author jerry 2022-09-24 06:41:45
// @Param client
// @Param url
// @Param reqBody
// @Return []byte
func HttpPostWithClient(client *http.Client, url string, reqBody []byte) []byte {
	//client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
	buffer := bytes.NewBuffer(reqBody)
	resp, err := client.Post(url, jsonContentType, buffer)
	//certInfo := resp.TLS.PeerCertificates[0]
	//fmt.Println("过期时间:", certInfo.NotAfter)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return respBody
}
