package server

import (
	"SGX_blockchain/src/accounts"
	"SGX_blockchain/src/utils"
	"encoding/hex"
	"encoding/json"
	"io"
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
		return nil, false, "Data count doesn't match!"
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

func (m *MainHandler) AccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequest(body, 2, 2)
		if valid {
			addressCompressed := crypto.ToCompressedPubKey(address)
			val, exists := m.d.Get(addressCompressed)
			//fmt.Println(utils.EncodeBytesToHexStringWith0x(addressCompressed))
			//fmt.Println(val, exists)
			var account *accounts.ExternalAccount
			if !exists {
				account = accounts.NewAccount()
				account.Id = addressCompressed
			} else {
				account = &accounts.ExternalAccount{}
				account.UnmarshalMsg(val)
			}

			account.StoreContract()

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
		//fmt.Println(utils.EncodeBytesToHexStringWith0x(addressCompressed))
		//fmt.Println(val, exists)
		var account = accounts.NewAccount()
		if val, exists := m.d.Get(addressCompressed); !exists {
			account.Id = addressCompressed
		} else {
			account.UnmarshalMsg(val)
		}
		code := gjson.GetBytes(body, "data.code").String()
		codeHash := gjson.GetBytes(body, "data.codeHash").String()
		if codeHash != hex.EncodeToString(crypto.Keccak256([]byte(code))) {
			w.WriteHeader(http.StatusOK)
			//TODO! 补充状态
		}

		codeHashBytes, _ := hex.DecodeString(codeHash)
		ok, txHash := account.StoreContract(codeHashBytes)

		if ok != nil {

		}

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

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
