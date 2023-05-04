package server

import (
	"SGX_blockchain/src/accounts"
	"SGX_blockchain/src/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"strconv"

	//"fmt"
	"github.com/tidwall/sjson"
	"net/http"
	"time"

	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/db"

	"github.com/tidwall/gjson"
)

const WrongResponse string = `{"status":"wrong"}`
const SuccessResponse string = `{"status":"ok"}`

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

type ContractDeployRequest struct {
	Data struct {
		From     string `json:"from"`
		Code     string `json:"code"`
		CodeHash string `json:"codeHash"`
		Ts       int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type ContractDeployResponse struct {
	Status      string `json:"status"`
	Transaction struct {
		From  string `json:"from"`
		Hash  string `json:"hash"`
		Nonce int64  `json:"nonce"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

type ContractCallRequest struct {
	Data struct {
		CodeHash     string `json:"codeHash"`
		From         string `json:"from"`
		FunctionName string `json:"functionName"`
		Params       string `json:"params"`
		Ts           int64  `json:"ts"`
	} `json:"data"`
	Signature string `json:"signature"`
}

type ContractCallResponse struct {
	Status      string `json:"status"`
	Transaction struct {
		CodeHash string `json:"codeHash"`
		From     string `json:"from"`
		Hash     string `json:"hash"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

//type WrongResponseStruct struct {
//	status string
//	info string
//}

const jsonString = `{"name":"last"}`

// func verify
type MainHandler struct {
	d db.Database
}

func NewHandler(d db.Database) *MainHandler {
	return &MainHandler{d}
}

// CheckRequest
// @Description 检查request body的合法性：1.格式；2.签名；
// @Author jerry 2022-09-24 06:46:57 ${time}
// @Param requestBody
// @Param count 一级字段下key的数量
// @Param dataCount data字段下key的数量
// @Return bool
func CheckRequest(requestBody []byte, count, dataCount int64) (address []byte, valid bool, errString string) {
	if !gjson.ValidBytes(requestBody) {
		return nil, false, "Invalid json syntax!"
	}
	//是否是json格式

	if gjson.GetBytes(requestBody, "@keys.#").Int() != count {
		return nil, false, "Key count doesn't match!"
	}
	//一级字段下key的数量

	data := gjson.GetBytes(requestBody, "data")
	if !data.Exists() || gjson.GetBytes(requestBody, "data.@keys.#").Int() != dataCount {
		return nil, false, "Data count doesn't match! " + gjson.GetBytes(requestBody, "data.@keys.#").String() + "!=" + strconv.Itoa(int(dataCount))
	}
	//data字段下key的数量
	dataMap := data.Map()
	if value, exists := dataMap["address"]; exists {
		address, _ = hex.DecodeString(value.String()[2:])
	} else {
		address, _ = hex.DecodeString(dataMap["from"].String()[2:])
	}
	sig := gjson.GetBytes(requestBody, "signature")

	if len(address) != 65 && len(address) != 33 {
		return nil, false, "Address Length doesn't match!"
	}
	if !sig.Exists() {
		return nil, false, "No signature!"
	}
	signature, _ := hex.DecodeString(sig.String()[2:])
	sigDER := crypto.EthereumSignatureToDER(signature)
	return address, crypto.VerifyMessageSignature(sigDER, address, []byte(data.String())), "Signature verification failed"
	//  若验证成功，则string返回值无意义，因此直接返回失败
}

func (m *MainHandler) GetAccount(addressCompressed []byte) (*accounts.ExternalAccount, bool) {
	val, exists := m.d.Get(addressCompressed)
	if !exists {
		return accounts.NewAccount(), false
	} else {
		account := &accounts.ExternalAccount{}
		account.UnmarshalMsg(val)
		return account, true
	}
}

func (m *MainHandler) AccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequest(body, 2, 2)
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
			r, _ := json.Marshal(resp)
			w.Write(r)
		} else {
			w.WriteHeader(http.StatusForbidden)
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
		address, valid, errString := CheckRequest(body, 2, 5)
		if valid {
			addressCompressed := crypto.ToCompressedPubKey(address)
			account, _ := m.GetAccount(addressCompressed)

			fileHashRaw := gjson.GetBytes(body, "data.fileHash")
			fileHashString := fileHashRaw.String()
			fileHashBytes := utils.DecodeHexStringToBytesWith0x(fileHashString)

			_, txHash, nonce := account.StoreFile(fileHashBytes)

			Content := gjson.GetBytes(body, "data.content")
			contentBytes, _ := base64.StdEncoding.DecodeString(Content.String())

			m.d.StoreFile(utils.EncodeBytesToHexStringWith0x(addressCompressed)+fileHashString, contentBytes)

			resp := &FileStoreResponse{
				Status: "ok",
				Transaction: struct {
					Hash     string `json:"hash"`
					FileHash string `json:"fileHash"`
					From     string `json:"from"`
					Nonce    int64  `json:"nonce"`
				}{
					Hash:     utils.EncodeBytesToHexStringWith0x(txHash),
					FileHash: fileHashString,
					From:     utils.EncodeBytesToHexStringWith0x(address),
					Nonce:    nonce},
				Ts: time.Now().UnixMilli(),
			}
			w.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(resp)
			w.Write(r)
		} else {
			w.WriteHeader(http.StatusForbidden)
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
		address, valid, errString := CheckRequest(body, 2, 3)
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
			w.WriteHeader(http.StatusForbidden)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) ContractDeployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequest(body, 2, 2)
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
		}

		addressCompressed := crypto.ToCompressedPubKey(address)

		var account = accounts.NewAccount()
		if val, exists := m.d.Get(addressCompressed); !exists {
			account.Id = addressCompressed
		} else {
			account.UnmarshalMsg(val)
		}
		code := gjson.GetBytes(body, "data.code").String()
		codeByte := []byte(code)
		codeHash := gjson.GetBytes(body, "data.codeHash").String()
		if codeHash != hex.EncodeToString(crypto.Keccak256(codeByte)) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(WrongResponse))
			//TODO! 补充状态
		}

		codeHashBytes, _ := hex.DecodeString(codeHash)
		ok, txHash, nonce := account.StoreContract(codeHashBytes)
		m.d.StoreContract(codeHashBytes, codeByte)
		if !ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(WrongResponse))
		}

		resp := &ContractDeployResponse{
			Status: "ok",
			Transaction: struct {
				From  string `json:"from"`
				Hash  string `json:"hash"`
				Nonce int64  `json:"nonce"`
			}{
				From:  utils.EncodeBytesToHexStringWith0x(address),
				Hash:  utils.EncodeBytesToHexStringWith0x(txHash),
				Nonce: nonce,
			},
			Ts: time.Now().UnixMilli(),
		}

		w.WriteHeader(http.StatusOK)
		r, _ := json.Marshal(resp)
		w.Write(r)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
