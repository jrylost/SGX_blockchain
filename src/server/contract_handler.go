package server

import (
	"SGX_blockchain/src/accounts"
	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/utils"
	"SGX_blockchain/src/vm"
	"SGX_blockchain/src/vm/ContractContext"
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
	Status string `json:"status"`
	Data   struct {
		From            string `json:"from"`
		Hash            string `json:"hash"`
		ContractAddress string `json:"contractAddress"`
		Nonce           int64  `json:"nonce"`
		CodeHash        string `json:"codeHash"`
	} `json:"data"`
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
	Status string `json:"status"`
	Data   struct {
		Hash            string `json:"hash"`
		Result          string `json:"result"`
		From            string `json:"from"`
		ContractAddress string `json:"contractAddress"`
		Nonce           int64  `json:"nonce"`
	} `json:"data"`
	Ts int64 `json:"ts"`
}

func (m *MainHandler) ContractDeployHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		address, valid, errString := CheckRequestWithAddress(body, 2, 6)
		if !valid {
			w.WriteHeader(http.StatusOK)
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
			errorString, _ := sjson.Set(WrongResponse, "error", `incorrect codeHash`)
			w.Write([]byte(errorString))
			return
		}
		contractABI := gjson.GetBytes(body, "data.abi").String()
		abiParsed, err := ContractContext.ABIParser(contractABI)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			errString, _ = sjson.Set(WrongResponse, "error", err.Error())
			w.Write([]byte(errString))
			return
		}
		codeHashBytes := utils.DecodeHexStringToBytesWith0x(codeHash)
		ok, txHash, nonce, contractAddress := account.StoreContract(codeHashBytes)
		if !ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(WrongResponse))
			return
		}
		m.d.StoreContract(utils.JoinBytes(addressCompressed, codeHashBytes), codeByte, abiParsed)
		ctime := time.Now().UnixMilli()
		txHashWith0x := utils.EncodeBytesToHexStringWith0x(txHash)

		resp := &ContractDeployResponse{
			Status: "ok",
			Data: struct {
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
		m.d.CreateContext(utils.JoinBytes(addressCompressed, codeHashBytes))
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
	//defer func() {
	//	// recover from panic if one occured. Set err to nil otherwise.
	//	if err := recover(); err != nil {
	//		w.WriteHeader(http.StatusOK)
	//		errorString, _ := sjson.Set(WrongResponse, "error", err)
	//		w.Write([]byte(errorString))
	//		return
	//	}
	//}()
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

		_, txHash := account.CallContract(contractAddress)
		txHashWith0x := utils.EncodeBytesToHexStringWith0x(txHash)

		codeHash := gjson.GetBytes(body, "data.codeHash").String()
		codeHashByte := utils.DecodeHexStringToBytesWith0x(codeHash)

		if contractAddress != utils.EncodeBytesToHexStringWith0x(crypto.Keccak256(account.Id, codeHashByte)) {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", `contractAddress and codeHash mismatch`)
			w.Write([]byte(errorString))
			return
		}

		contractContent, ok := m.d.GetContract(utils.JoinBytes(addressCompressed, codeHashByte))
		if !ok {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", "not deployed")
			w.Write([]byte(errorString))
			return
		}
		ctx, abi, ok := m.d.GetContext(utils.JoinBytes(addressCompressed, codeHashByte))
		if !ok {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", "not deployed")
			w.Write([]byte(errorString))
			return
		}

		functionName := gjson.GetBytes(body, "data.functionName").String()

		functionInputs := gjson.GetBytes(body, "data.functionInputs").String()
		values, err := ContractContext.ContractInputHandler(functionInputs)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", err.Error())
			w.Write([]byte(errorString))
			return
		}

		var abiFunction ContractContext.ContractFunction
		var flag = false
		for _, function := range abi.ContractFunctions {
			if function.FunctionName == functionName {
				abiFunction = function
				flag = true
				break
			}
		}
		if !flag {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", "incorrect functionName")
			w.Write([]byte(errorString))
			return
		}
		inputs, err := ContractContext.ContractInputVerify(values, abiFunction.FunctionInputs)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", err.Error())
			w.Write([]byte(errorString))
			return
		}

		virtualMachine := vm.NewVirtualMachine(contractContent)
		results, newctx, err := virtualMachine.Call(string(addressCompressed), txHashWith0x, abi.ContractName, functionName, contractContent, inputs, ctx)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", err.Error())
			w.Write([]byte(errorString))
			return
		}
		jsonResults, err := json.Marshal(results)

		if err != nil {
			w.WriteHeader(http.StatusOK)
			errorString, _ := sjson.Set(WrongResponse, "error", err.Error())
			w.Write([]byte(errorString))
			return
		}
		m.d.StoreContext(utils.JoinBytes(addressCompressed, codeHashByte), newctx)
		//m.d.StoreContract(contractAddress, codeByte, contractABI)
		//if !ok {
		//	w.WriteHeader(http.StatusOK)
		//	w.Write([]byte(WrongResponse))
		//}
		ctime := time.Now().UnixMilli()

		resp := &ContractCallResponse{
			Status: "ok",
			Data: struct {
				Hash            string `json:"hash"`
				Result          string `json:"result"`
				From            string `json:"from"`
				ContractAddress string `json:"contractAddress"`
				Nonce           int64  `json:"nonce"`
			}{
				Hash:            txHashWith0x,
				Result:          string(jsonResults),
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
