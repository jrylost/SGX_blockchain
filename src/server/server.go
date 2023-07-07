package server

import (
	"SGX_blockchain/src/accounts"
	"SGX_blockchain/src/utils"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"sync"

	//"encoding/json"
	"io"
	"strconv"

	//"fmt"
	"github.com/tidwall/sjson"
	"net/http"
	"time"

	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/db"

	"github.com/tidwall/gjson"

	jsoniter "github.com/json-iterator/go"
)

const WrongResponse string = `{"status":"error"}`
const JsonMarshalErrorResponse string = `{"status":"error","error":"json marshal error"}`

//const SuccessResponse string = `{"status":"ok"}`

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var mu sync.Mutex

type BlockInfoRequest struct {
	Data struct {
		Number int64 `json:"number"`
		Ts     int64 `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type BlockInfoResponse struct {
	Status string `json:"status"`
	Data   struct {
		Number       int64    `json:"number"`
		Transactions []string `json:"transactions"`
		Count        int64    `json:"count"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type TransactionInfoRequest struct {
	Data struct {
		Hash string `json:"hash"`
		Ts   int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type TransactionInfoResponse struct {
	Status string `json:"status"`
	Data   struct {
		Hash          string `json:"hash"`
		Type          string `json:"type"`
		TransactionTs int64  `json:"transactionTs"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type AccountInfoRequest struct {
	Data struct {
		Address string `json:"address"`
		Ts      int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type AccountInfoResponse struct {
	Status string `json:"status"`
	Data   struct {
		Address string `json:"address"`
		Balance int64  `json:"balance"`
		Nonce   int64  `json:"nonce"`
		Count   int64  `json:"count"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

type TransferRequest struct {
	Data struct {
		From  string `json:"from"`
		To    string `json:"to"`
		Nonce int64  `json:"nonce"`
		Value int64  `json:"value"`
		Ts    int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type TransferResponse struct {
	Status      string `json:"status"`
	Transaction struct {
		Hash  string `json:"hash"`
		From  string `json:"from"`
		To    string `json:"to"`
		Nonce int64  `json:"nonce"`
		Value int64  `json:"value"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

type FileStoreRequest struct {
	Data struct {
		Content  string `json:"content"`
		FileHash string `json:"fileHash"`
		From     string `json:"from"`
		Nonce    int64  `json:"nonce"`
		Ts       int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type FileStoreResponse struct {
	Status      string `json:"status"`
	Transaction struct {
		Hash     string `json:"hash"`
		FileHash string `json:"fileHash"`
		From     string `json:"from"`
		Nonce    int64  `json:"nonce"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

type FileRetrieveRequest struct {
	Data struct {
		From     string `json:"from"`
		FileHash string `json:"fileHash"`
		Ts       int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type FileRetrieveResponse struct {
	Status      string `json:"status"`
	Transaction struct {
		FileHash string `json:"fileHash"`
		From     string `json:"from"`
		Content  string `json:"content"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

type KVStoreRequest struct {
	Data struct {
		From  string `json:"from"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Nonce int64  `json:"nonce"`
		Ts    int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type KVStoreResponse struct {
	Status      string `json:"status"`
	Transaction struct {
		Hash  string `json:"hash"`
		From  string `json:"from"`
		Key   string `json:"key"`
		Nonce int64  `json:"nonce"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

type KVRetrieveRequest struct {
	Data struct {
		From string `json:"from"`
		Key  string `json:"key"`
		Ts   int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type KVRetrieveResponse struct {
	Status      string `json:"status"`
	Transaction struct {
		From  string `json:"from"`
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

//type WrongResponseStruct struct {
//	status string
//	info string
//}

//const jsonString = `{"name":"last"}`

// func verify
type MainHandler struct {
	d db.Database
}

func NewHandler(d db.Database) *MainHandler {
	return &MainHandler{d}
}

// CheckRequestWithAddress
// @Description 检查request body的合法性：1.格式；2.签名；
// @Author jerry 2022-09-24 06:46:57 ${time}
// @Param requestBody
// @Param count 一级字段下key的数量
// @Param dataCount data字段下key的数量
// @Return bool
func CheckRequestWithAddress(requestBody []byte, count, dataCount int64) ([]byte, bool, string) {
	if !gjson.ValidBytes(requestBody) {
		return nil, false, "invalid json syntax"
	}
	//是否是json格式

	if gjson.GetBytes(requestBody, "@keys.#").Int() != count {
		return nil, false, "key count error"
	}
	//一级字段下key的数量

	data := gjson.GetBytes(requestBody, "data")
	if !data.Exists() || gjson.GetBytes(requestBody, "data.@keys.#").Int() != dataCount {
		return nil, false, "data count error " + gjson.GetBytes(requestBody, "data.@keys.#").String() + "!=" + strconv.Itoa(int(dataCount))
	}
	//data字段下key的数量
	dataMap := data.Map()
	var address []byte
	var err error
	if value, exists := dataMap["address"]; exists {
		address, err = hex.DecodeString(value.String()[2:])
	} else {
		address, err = hex.DecodeString(dataMap["from"].String()[2:])
	}
	if err != nil {
		return nil, false, "hex address decoding error"
	}
	sig := gjson.GetBytes(requestBody, "signature")

	if len(address) != 65 && len(address) != 33 {
		return nil, false, "address length error"
	}
	if !sig.Exists() {
		return nil, false, "no signature"
	}
	signature, _ := hex.DecodeString(sig.String()[2:])
	sigDER := crypto.EthereumSignatureToDER(signature)
	return address, crypto.VerifyMessageSignature(sigDER, address, []byte(data.String())), "signature verification error"
	//  若验证成功，则string返回值无意义，因此直接返回失败
}

func (m *MainHandler) GetAccount(addressCompressed []byte) (*accounts.ExternalAccount, bool) {
	val, exists := m.d.Get(addressCompressed)
	if !exists {
		return accounts.NewAccount(), false
	} else {
		account := &accounts.ExternalAccount{}
		_, err := account.UnmarshalMsg(val)
		if err != nil {
			return nil, false
		}
		return account, true
	}
}

//func (m *MainHandler) SaveAccount(addressCompressed []byte, *accounts.ExternalAccount) bool {
//	val, exists := m.d.Get(addressCompressed)
//	m.d.
//
//	if !exists {
//		return  false
//	} else {
//		account := &accounts.ExternalAccount{}
//		_, err := account.UnmarshalMsg(val)
//		if err != nil {
//			return nil, false
//		}
//		return account, true
//	}
//}

func (m *MainHandler) BlockInfoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
	}

	defer r.Body.Close()
	switch r.Method {
	case "POST":
		var valid = true
		var errString string
		//address, valid, errString := CheckRequestWithAddress(body, 2, 2)
		if gjson.GetBytes(body, "@keys.#").Int() != 2 {
			valid = false
			errString = "key count error"
		}
		//一级字段下key的数量

		data := gjson.GetBytes(body, "data")
		if !data.Exists() || gjson.GetBytes(body, "data.@keys.#").Int() != 2 {
			valid = false
			errString = "data count error"
		}

		//valid := true
		if valid {
			//addressCompressed := crypto.ToCompressedPubKey(address)
			//account, _ := m.GetAccount(addressCompressed)
			number := gjson.GetBytes(body, "data.number").Int()
			txs := m.d.GetTxFromBlock(number)
			resp := &BlockInfoResponse{
				Status: "ok",
				Data: struct {
					Number       int64    `json:"number"`
					Transactions []string `json:"transactions"`
					Count        int64    `json:"count"`
				}{
					Number:       number,
					Transactions: txs,
					Count:        int64(len(txs)),
				},
				Ts: time.Now().UnixMilli(),
			}
			r, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, errWrite := w.Write([]byte(JsonMarshalErrorResponse))
				if errWrite != nil {
					log.Println(errWrite.Error())
				}
				return
			}
			w.WriteHeader(http.StatusOK)

			_, errWrite := w.Write(r)
			if errWrite != nil {
				log.Println(errWrite.Error())
			}
		} else {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			_, errWrite := w.Write([]byte(errorString))
			if errWrite != nil {
				log.Println(errWrite.Error())
			}
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) TransactionInfoHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		var valid = true
		var errString string
		//address, valid, errString := CheckRequestWithAddress(body, 2, 2)
		if gjson.GetBytes(body, "@keys.#").Int() != 2 {
			valid = false
			errString = "key count error"
		}
		//一级字段下key的数量

		data := gjson.GetBytes(body, "data")
		if !data.Exists() || gjson.GetBytes(body, "data.@keys.#").Int() != 2 {
			valid = false
			errString = "data count error"
		}

		//valid := true
		if valid {
			//addressCompressed := crypto.ToCompressedPubKey(address)
			//account, _ := m.GetAccount(addressCompressed)
			txhash := gjson.GetBytes(body, "data.hash").String()
			txs := m.d.GetTx(txhash)
			resp, _ := sjson.SetBytes([]byte(txs), "ts", time.Now().UnixMilli())

			w.WriteHeader(http.StatusOK)
			//result, err := json.Marshal(resp)
			//if err != nil {
			//	return
			//}

			w.Write(resp)
		} else {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) AccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	var jsonlib = jsoniter.ConfigCompatibleWithStandardLibrary

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 2)
		if valid {
			addressCompressed := crypto.ToCompressedPubKey(address)
			account, _ := m.GetAccount(addressCompressed)

			resp := &AccountInfoResponse{
				Status: "ok",
				Data: struct {
					Address string `json:"address"`
					Balance int64  `json:"balance"`
					Nonce   int64  `json:"nonce"`
					Count   int64  `json:"count"`
				}{Address: utils.EncodeBytesToHexStringWith0x(address),
					Balance: account.Balance,
					Nonce:   account.Nonce,
					Count:   account.Count},
				Ts: time.Now().UnixMilli(),
			}
			w.WriteHeader(http.StatusOK)
			r, err := jsonlib.Marshal(resp)
			if err != nil {
				return
			}

			w.Write(r)
		} else {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) FileStoreHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 5)
		if valid {
			addressCompressed := crypto.ToCompressedPubKey(address)
			account, _ := m.GetAccount(addressCompressed)

			fileHashRaw := gjson.GetBytes(body, "data.fileHash")
			if !fileHashRaw.Exists() {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"wrong","error":"fileHash missing"}`))
				return
			}
			fileHashString := fileHashRaw.String()
			fileHashBytes := utils.DecodeHexStringToBytesWith0x(fileHashString)

			_, txHash, nonce := account.StoreFile(fileHashBytes)

			Content := gjson.GetBytes(body, "data.content")
			contentBytes, _ := base64.StdEncoding.DecodeString(Content.String())

			m.d.StoreFile(utils.EncodeBytesToHexStringWith0x(addressCompressed)+fileHashString, contentBytes)
			ctime := time.Now().UnixMilli()
			txHashWith0x := utils.EncodeBytesToHexStringWith0x(txHash)

			resp := &FileStoreResponse{
				Status: "ok",
				Transaction: struct {
					Hash     string `json:"hash"`
					FileHash string `json:"fileHash"`
					From     string `json:"from"`
					Nonce    int64  `json:"nonce"`
				}{
					Hash:     txHashWith0x,
					FileHash: fileHashString,
					From:     utils.EncodeBytesToHexStringWith0x(address),
					Nonce:    nonce},
				Ts: ctime,
			}
			fmt.Println("Nonce here:", account.Nonce)
			w.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(resp)
			m.d.StoreTxToBlock(ctime, txHashWith0x)
			m.d.StoreTx(txHashWith0x, "File store", ctime)
			//accountjsonrefreshed, _ := account.MarshalMsg(nil)
			accountjsonrefreshed, err := account.MarshalMsg(nil)
			if err != nil {
				fmt.Println(err.Error())
			}
			mu.Lock()
			m.d.Put(addressCompressed, accountjsonrefreshed)
			mu.Unlock()
			w.Write(r)
		} else {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) FileRetrieveHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 3)
		if valid {
			addressCompressed := crypto.ToCompressedPubKey(address)
			//account, _ := m.GetAccount(addressCompressed)

			fileHashRaw := gjson.GetBytes(body, "data.fileHash")
			fileHashString := fileHashRaw.String()

			contentBytes := m.d.RetrieveFile(utils.EncodeBytesToHexStringWith0x(addressCompressed) + fileHashString)
			content := base64.StdEncoding.EncodeToString(contentBytes)

			resp := &FileRetrieveResponse{
				Status: "ok",
				Transaction: struct {
					FileHash string `json:"fileHash"`
					From     string `json:"from"`
					Content  string `json:"content"`
				}{
					FileHash: fileHashString,
					From:     utils.EncodeBytesToHexStringWith0x(address),
					Content:  content,
				},
				Ts: time.Now().UnixMilli(),
			}

			w.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(resp)
			w.Write(r)
		} else {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) KVStoreHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 5)
		if valid {
			addressCompressed := crypto.ToCompressedPubKey(address)
			account, _ := m.GetAccount(addressCompressed)

			keyRaw := gjson.GetBytes(body, "data.key")
			keyString := keyRaw.String()
			keyHashBytes := crypto.Keccak256([]byte(keyString))

			_, txHash, nonce := account.StoreKV([]byte(keyString))

			Value := gjson.GetBytes(body, "data.value")
			valueBytes := []byte(Value.String())

			m.d.StoreKV(utils.EncodeBytesToHexStringWith0x(utils.JoinBytes(addressCompressed, keyHashBytes)), valueBytes)
			ctime := time.Now().UnixMilli()
			txHashWith0x := utils.EncodeBytesToHexStringWith0x(txHash)

			resp := &KVStoreResponse{
				Status: "ok",
				Transaction: struct {
					Hash  string `json:"hash"`
					From  string `json:"from"`
					Key   string `json:"key"`
					Nonce int64  `json:"nonce"`
				}{
					Hash:  txHashWith0x,
					From:  utils.EncodeBytesToHexStringWith0x(address),
					Key:   keyString,
					Nonce: nonce,
				},
				Ts: ctime,
			}

			w.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(resp)

			m.d.StoreTxToBlock(ctime, txHashWith0x)
			m.d.StoreTx(txHashWith0x, "KV store", ctime)
			accountjsonrefreshed, err := account.MarshalMsg(nil)
			if err != nil {
				fmt.Println(err.Error())
			}
			mu.Lock()
			m.d.Put(addressCompressed, accountjsonrefreshed)
			mu.Unlock()
			//accountjsonrefreshed, _ := json.Marshal(account)
			w.Write(r)
		} else {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) KVRetrieveHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 3)
		if valid {
			addressCompressed := crypto.ToCompressedPubKey(address)
			//account, _ := m.GetAccount(addressCompressed)

			keyRaw := gjson.GetBytes(body, "data.key")
			keyString := keyRaw.String()
			keyHashBytes := crypto.Keccak256([]byte(keyString))

			//_, txHash, nonce := account.StoreKV([]byte(keyString))

			valueBytes := m.d.RetrieveKV(utils.EncodeBytesToHexStringWith0x(bytes.Join([][]byte{addressCompressed, keyHashBytes}, []byte(""))))

			resp := &KVRetrieveResponse{
				Status: "ok",
				Transaction: struct {
					From  string `json:"from"`
					Key   string `json:"key"`
					Value string `json:"value"`
				}{
					From:  utils.EncodeBytesToHexStringWith0x(address),
					Key:   keyString,
					Value: string(valueBytes),
				},
				Ts: time.Now().UnixMilli(),
			}

			w.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(resp)
			w.Write(r)
		} else {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
