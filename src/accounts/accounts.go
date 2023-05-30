package accounts

import (
	"SGX_blockchain/src/crypto"
	"encoding/binary"
	"encoding/hex"
)

// Account Account接口，要求实现转账transfer函数
type Account interface {
	Transfer(*ExternalAccount, int64) bool
	Serialize() []byte
	Deserialize([]byte)
	StoreFile(fileHash []byte)
}

//go:generate msgp

type ExternalAccount struct {
	Version  int64    `msg:"version"`
	Id       []byte   `msg:"id"`       //帐号id为33字节
	Count    int64    `msg:"count"`    //文件数量
	Nonce    int64    `msg:"nonce"`    //计数器，确保交易只执行一次
	Balance  int64    `msg:"balance"`  //余额
	File     []string `msg:"file"`     //文件
	Key      []string `msg:"key"`      //key-value存储
	Contract []string `msg:"contract"` //智能合约
}

// CreateAccount 传入id以及余额，创建一个新的帐号
func CreateAccount(id []byte, balance int64) *ExternalAccount {
	return &ExternalAccount{
		Id:      id,
		Nonce:   0,
		Balance: balance,
		//File:    make([]string, 0),
	}
}

// NewAccount 传入id以及余额，创建一个新的帐号
func NewAccount() *ExternalAccount {
	return &ExternalAccount{
		//File: make([]string),
	}
}

// Transfer ExternalAccount给recipient转账，金额为amount
func (account *ExternalAccount) Transfer(recipient *ExternalAccount, amount int64) (bool, []byte) {
	if account.Balance < amount {
		return false, []byte("Insufficient funds!")
	}
	//随机数自增，留用作同步确认
	account.Nonce++
	recipient.Nonce++
	//transferor减去对应金额，recipient增加对应金额
	account.Balance -= amount
	recipient.Balance += amount

	accountNonceByte := make([]byte, 8)
	recipientNonceByte := make([]byte, 8)
	binary.PutVarint(accountNonceByte, account.Nonce)
	binary.PutVarint(recipientNonceByte, recipient.Nonce)
	return true, crypto.Keccak256(account.Id, recipient.Id, accountNonceByte, recipientNonceByte)
}

func (account *ExternalAccount) StoreFile(fileHash []byte) (bool, []byte, int64) {
	account.File = append(account.File, hex.EncodeToString(fileHash))
	account.Nonce++
	account.Count++
	accountNonceByte := make([]byte, 8)
	binary.PutVarint(accountNonceByte, account.Nonce)
	return true, crypto.Keccak256(account.Id, accountNonceByte, fileHash), account.Nonce
}

func (account *ExternalAccount) StoreContract(contractHash []byte) (bool, []byte, int64) {
	account.File = append(account.File, hex.EncodeToString(contractHash))
	account.Nonce++

	accountNonceByte := make([]byte, 8)
	binary.PutVarint(accountNonceByte, account.Nonce)
	return true, crypto.Keccak256(account.Id, accountNonceByte, contractHash), account.Nonce
}

func (account *ExternalAccount) StoreKV(key []byte) (bool, []byte, int64) {
	account.Key = append(account.Key, string(key))
	account.Nonce++
	accountNonceByte := make([]byte, 8)
	binary.PutVarint(accountNonceByte, account.Nonce)
	return true, crypto.Keccak256(account.Id, accountNonceByte, key), account.Nonce
}
