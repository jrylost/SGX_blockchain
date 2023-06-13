package server

import (
	"SGX_blockchain/src/accounts"
	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/utils"
	"encoding/hex"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	"net/http"
	"time"
)

type ContractDeployRequest struct {
	Data struct {
		Name     string `json:"name"`
		ABI      string `json:"abi"`
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
		From            string `json:"from"`
		Hash            string `json:"hash"`
		ContractAddress string `json:"contractAddress"`
		Nonce           int64  `json:"nonce"`
		CodeHash        string `json:"codeHash"`
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
		Hash            string `json:"hash"`
		Result          string `json:"result"`
		From            string `json:"from"`
		ContractAddress string `json:"contractAddress"`
		Nonce           int64  `json:"nonce"`
	} `json:"transaction"`
	Ts int64 `json:"ts"`
}

func (m *MainHandler) ContractDeployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 6)
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
			return
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
		if codeHash[2:] != hex.EncodeToString(crypto.Keccak256(codeByte)) {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", `code and codeHash mismatch`)
			w.Write([]byte(errorString))
		}
		contractABI := gjson.GetBytes(body, "data.abi").String()
		codeHashBytes := utils.DecodeHexStringToBytesWith0x(codeHash)
		ok, txHash, nonce, contractAddress := account.StoreContract(codeHashBytes)
		m.d.StoreContract(contractAddress, codeByte, contractABI)
		if !ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(WrongResponse))
		}
		ctime := time.Now().UnixMilli()
		txHashWith0x := utils.EncodeBytesToHexStringWith0x(txHash)

		resp := &ContractDeployResponse{
			Status: "ok",
			Transaction: struct {
				From            string `json:"from"`
				Hash            string `json:"hash"`
				ContractAddress string `json:"contractAddress"`
				Nonce           int64  `json:"nonce"`
				CodeHash        string `json:"codeHash"`
			}{From: utils.EncodeBytesToHexStringWith0x(address),
				Hash:            txHashWith0x,
				ContractAddress: utils.EncodeBytesToHexStringWith0x(contractAddress),
				Nonce:           nonce,
				CodeHash:        codeHash},

			Ts: ctime,
		}

		w.WriteHeader(http.StatusOK)
		r, _ := json.Marshal(resp)

		m.d.StoreTxToBlock(ctime, txHashWith0x)
		m.d.StoreTx(txHashWith0x, "Contract depoly", ctime)
		accountjsonrefreshed, err := account.MarshalMsg(nil)
		if err != nil {
			fmt.Println(err.Error())
		}
		mu.Lock()
		m.d.Put(addressCompressed, accountjsonrefreshed)
		mu.Unlock()

		w.Write(r)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *MainHandler) ContractCallHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 6)
		if !valid {
			w.WriteHeader(http.StatusForbidden)
			errorString, _ := sjson.Set(WrongResponse, "error", errString)
			w.Write([]byte(errorString))
			return
		}
		addressCompressed := crypto.ToCompressedPubKey(address)

		var account = accounts.NewAccount()
		if val, exists := m.d.Get(addressCompressed); !exists {
			account.Id = addressCompressed
		} else {
			account.UnmarshalMsg(val)
		}

		contractAddress := gjson.GetBytes(body, "data.contractAddress").String()
		codeHash := gjson.GetBytes(body, "data.codeHash").String()
		codeHashByte := utils.DecodeHexStringToBytesWith0x(codeHash)

		if contractAddress != utils.EncodeBytesToHexStringWith0x(crypto.Keccak256(account.Id, codeHashByte)) {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", `contractAddress and codeHash mismatch`)
			w.Write([]byte(errorString))
		}
		//m.d.StoreContract(contractAddress, codeByte, contractABI)
		//if !ok {
		//	w.WriteHeader(http.StatusOK)
		//	w.Write([]byte(WrongResponse))
		//}
		ctime := time.Now().UnixMilli()
		_, txHash := account.CallContract(contractAddress)
		txHashWith0x := utils.EncodeBytesToHexStringWith0x(txHash)

		resp := &ContractCallResponse{
			Status: "",
			Transaction: struct {
				Hash            string `json:"hash"`
				Result          string `json:"result"`
				From            string `json:"from"`
				ContractAddress string `json:"contractAddress"`
				Nonce           int64  `json:"nonce"`
			}{
				Hash:            txHashWith0x,
				Result:          "",
				From:            utils.EncodeBytesToHexStringWith0x(address),
				ContractAddress: contractAddress,
				Nonce:           account.Nonce,
			},
			Ts: ctime,
		}
		w.WriteHeader(http.StatusOK)
		r, _ := json.Marshal(resp)

		m.d.StoreTxToBlock(ctime, txHashWith0x)
		m.d.StoreTx(txHashWith0x, "Contract Call", ctime)
		accountjsonrefreshed, err := account.MarshalMsg(nil)
		if err != nil {
			fmt.Println(err.Error())
		}
		mu.Lock()
		m.d.Put(addressCompressed, accountjsonrefreshed)
		mu.Unlock()

		w.Write(r)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
