package client

import (
	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/server"
	"SGX_blockchain/src/utils"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"
)

var Response = []string{
	`获取远程认证：`,
	`{
    "status":"ok",
    "data":{
        "report":"eX26q+uozgOQvTmTkaUj1ABv7oIdZjMQ1h4p51fi7Di7rGD+tlZG9SNfkRAVujoOohkg8laRsTAtQiSChPj4/VfrHMDVRp/bQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896P+CHQ4FLbCm5UOdJxGcoLFUAJUfjsAm+ZYvW6hnvcRpjdq4hqrp74E5+fuvS35O1Mxni9RRt9zkcBDE2b5QS0GCFTHw+zSrnF9AU8RxfzIhEpEu9u/7iw2d8eqEMZotf5CnR3ypt7mlErkz8nMAnUG6CSjzpT78LnrPZALp4L3LuuuPhmstAJB+MK2dxD8J0wEG8qc8ZtEHLd3DfsdKla6UdvSKDH70mXdXJKLNjtGwOUY/rPzheAwBAcXoD60o42RrmMQCOw2z99zZHGOe2GIp3jypF7XkUW3tXsgwBMbjcnP0+yAoGy0tzC7oQushNfMYhfw/EKVW1PvoD58OKTaGotc9/tDVrE1XcJAMA4/iNbDWlngscX/d3PHYwSLcdP51cb+nJKJ4tACbQNZCrwPrLKNleiQ6/84mUUPNOzXzR1ShC5VTeZl9SyIZqOfxSk81MkzAD0PnqWbVHhdjJUYjzy/S0cD/cVQdbzkHfHDKgQIMhBv0HxU24o5In+r/PbLd3MlQMjV/toYPylk8=",
        "cert":"DVRp/bQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896P+CHQ4FLbCm5UOdJxGcoLbQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896PbQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896PbQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896PbQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896P",
    }
}`, `查询账户信息：`,
	`{
    "data":{
        "address":"0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
        "ts":1650333610000,
    },
    "signature":"0x4c49d393b56749d6a2048f2ef6eaa60dba54b45d78f3d0ce9bccb97f6f1e884b"
}`, `返回账户信息：`,
	`{
    "status":"ok",
    "data":{
        "address":"0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
        "balance":10000000000,
        "nonce":10,
        "count":10,
    },
    "ts":1650333610000,
}`, `申请转账：`,
	`{
    "data":{
        "from":"0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
        "to":"0x1c1643c57b0ec7542498f693f201da4ccbb9991289e45631436bad366fdc111d",
        "nonce":15,
        "value":100000000,
        "ts":1650333610000,
    },
    "signature":"0x4c49d393b56749d6a2048f2ef6eaa60dba54b45d78f3d0ce9bccb97f6f1e884b"
}`, `返回转账结果：`,
	`{
    "status":"ok",
    "data":{
        "transaction":{
            "hash": "0xcc7c9dbe7bb4e409967803c6a2c4859e5068d4044ff7cf91a1c5179b92bbf967",
            "from": "0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
            "to": "0x1c1643c57b0ec7542498f693f201da4ccbb9991289e45631436bad366fdc111d",
            "nonce": 15,
            "value": 100000000,
        }
    },
    "ts":1650333610000,
}`, `申请存储文件：`,
	`{
    "data":{
        "address":"0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
        "content":"eX26q+uozgOQvTmTkaUj1ABv7oIdZjMQ1h4p51fi7Di7rGD+tlZG9SNfkRAVujoOohkg8laRsTAtQiSChPj4/VfrHMDVRp/bQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896P+CHQ4FLbCm5UOdJxGcoLFUAJUfjsAm+ZYvW6hnvcRpjdq4hqrp74E5+fuvS35O1Mxni9RRt9zkcBDE2b5QS0GCFTHw+zSrnF9AU8RxfzIhEpEu9u/7iw2d8eqEMZotf5CnR3ypt7mlErkz8nMAnUG6CSjzpT78LnrPZALp4L3LuuuPhmstAJB+MK2dxD8J0wEG8qc8ZtEHLd3DfsdKla6UdvSKDH70mXdXJKLNjtGwOUY/rPzheAwBAcXoD60o42RrmMQCOw2z99zZHGOe2GIp3jypF7XkUW3tXsgwBMbjcnP0+yAoGy0tzC7oQushNfMYhfw/EKVW1PvoD58OKTaGotc9/tDVrE1XcJAMA4/iNbDWlngscX/d3PHYwSLcdP51cb+nJKJ4tACbQNZCrwPrLKNleiQ6/84mUUPNOzXzR1ShC5VTeZl9SyIZqOfxSk81MkzAD0PnqWbVHhdjJUYjzy/S0cD/cVQdbzkHfHDKgQIMhBv0HxU24o5In+r/PbLd3MlQMjV/toYPylk8=",
        "filehash":"0xa4473b3f3a90025c936646d75195a3ab0a4685a31142423121375baed271dd6d",
        "nonce":15,
        "ts":1650333610000,
    },
    "signature":"0x4c49d393b56749d6a2048f2ef6eaa60dba54b45d78f3d0ce9bccb97f6f1e884b"
}`, `返回存储信息：`,
	`{
    "status":"ok",
    "data":{
        "transaction":{
            "hash": "0xcc7c9dbe7bb4e409967803c6a2c4859e5068d4044ff7cf91a1c5179b92bbf967",
            "address": "0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
            "filehash": "0x1c1643c57b0ec7542498f693f201da4ccbb9991289e45631436bad366fdc111d",
            "nonce": 15,
        }
    },
    "ts":1650333610000,
}`, `申请读取文件：`,
	`{
    "data":{
        "address":"0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
        "filehash":"0xa4473b3f3a90025c936646d75195a3ab0a4685a31142423121375baed271dd6d",
        "ts":1650333610000,
    },
    "signature":"0x4c49d393b56749d6a2048f2ef6eaa60dba54b45d78f3d0ce9bccb97f6f1e884b"
}`, `返回文件读取结果：`,
	`{
    "status":"ok",
    "data":{
        "filehash": "0xcc7c9dbe7bb4e409967803c6a2c4859e5068d4044ff7cf91a1c5179b92bbf967",
        "address": "0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
        "content":"eX26q+uozgOQvTmTkaUj1ABv7oIdZjMQ1h4p51fi7Di7rGD+tlZG9SNfkRAVujoOohkg8laRsTAtQiSChPj4/VfrHMDVRp/bQjvMyZ7xmDK1LTZAyIkDwOFnrhqqFr4gDUI9XS53bY6yLpnVFcN7e896P+CHQ4FLbCm5UOdJxGcoLFUAJUfjsAm+ZYvW6hnvcRpjdq4hqrp74E5+fuvS35O1Mxni9RRt9zkcBDE2b5QS0GCFTHw+zSrnF9AU8RxfzIhEpEu9u/7iw2d8eqEMZotf5CnR3ypt7mlErkz8nMAnUG6CSjzpT78LnrPZALp4L3LuuuPhmstAJB+MK2dxD8J0wEG8qc8ZtEHLd3DfsdKla6UdvSKDH70mXdXJKLNjtGwOUY/rPzheAwBAcXoD60o42RrmMQCOw2z99zZHGOe2GIp3jypF7XkUW3tXsgwBMbjcnP0+yAoGy0tzC7oQushNfMYhfw/EKVW1PvoD58OKTaGotc9/tDVrE1XcJAMA4/iNbDWlngscX/d3PHYwSLcdP51cb+nJKJ4tACbQNZCrwPrLKNleiQ6/84mUUPNOzXzR1ShC5VTeZl9SyIZqOfxSk81MkzAD0PnqWbVHhdjJUYjzy/S0cD/cVQdbzkHfHDKgQIMhBv0HxU24o5In+r/PbLd3MlQMjV/toYPylk8=",
    },
    "ts":1650333610000,
}`, `申请存储键值对：`,
	`{
    "data":{
        "address":"0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
        "key":"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "value":"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
        "nonce":15,
        "ts":1650333610000,
    },
    "signature":"0x4c49d393b56749d6a2048f2ef6eaa60dba54b45d78f3d0ce9bccb97f6f1e884b"
}`, `返回键值对存储结果：`,
	`{
    "status":"ok",
    "data":{
        "transaction":{
            "hash": "0xcc7c9dbe7bb4e409967803c6a2c4859e5068d4044ff7cf91a1c5179b92bbf967",
            "address": "0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
            "key": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
            "nonce": 15,
        }
    },
    "ts":1650333610000,
}`,
	`deployContract：`,
	`{
    "data":{
		"from": "0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
		"code":"",
		"codeHash": "",
    },
"signature":"0x4c49d393b56749d6a2048f2ef6eaa60dba54b45d78f3d0ce9bccb97f6f1e884b"
}`,
	`deployContract：`,
	`{
	"status":"ok",
    "transaction":{
		"from": "0x95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18",
		"hash":"",
    },
	"ts":1650333610000,
}`,
}

const ServerURL = "https://localhost:8888"
const jsonContentType string = "application/json"

type Client struct {
	tls  *tls.Config
	pair *crypto.KeyPair
}

func NewClient(t *tls.Config, pair *crypto.KeyPair) *Client {
	return &Client{tls: t, pair: pair}
}

func (c *Client) RequestForAccountInfo(pair *crypto.KeyPair) []byte {
	address := utils.EncodeBytesToHexStringWith0x(pair.PubK)
	data := struct {
		Address string `json:"address"`
		Ts      int64  `json:"ts"`
	}{
		address,
		time.Now().Unix(),
	}
	dataBytes, err := json.Marshal(data)
	//dataBytes, err := sjson.SetBytes(nil, "@this", data)
	if err != nil {
		fmt.Println("err json:", err)
	}
	//fmt.Println(string(dataBytes))
	sig := pair.SignMessage(dataBytes)
	//fmt.Println("data:", hex.EncodeToString(dataBytes))

	req := &server.AccountInfoRequest{
		Data:      data,
		Signature: utils.EncodeBytesToHexStringWith0x(sig),
	}
	reqBytes, _ := json.Marshal(req)
	//buffer := bytes.NewBuffer(reqBytes)
	resp := utils.HttpPost(c.tls, ServerURL+"/account/info", reqBytes)
	return resp
	//resp, err := http.Post(ServerURL+"/account/info", jsonContentType, buffer)
	//fmt.Println(resp)
	//fmt.Println(gjson.GetBytes(resp, "@pretty"))

}

//testserver := httptest.NewServer(http.HandlerFunc(h.AccountInfoHandler))
//accountsRequest := &AccountInfoRequest{
//Address:   utils.GenerateRandomHexStringWith0x(32),
//Timestamp: time.Now().Unix(),
//Signature: "",
//}
//resp, _ := http.Post(testserver.URL, jsonContentType, buffer)
//defer testserver.Close()
//bodyBytes, _ := json.Marshal(accountsRequest)
//buffer := bytes.NewBuffer(bodyBytes)
//
//defer resp.Body.Close()
