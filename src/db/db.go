package db

import (
	"SGX_blockchain/src/vm/ContractContext"
)

type Database interface {
	//Get returns the value for given key if it exists
	Get(s []byte) ([]byte, bool)
	Put(s, v []byte) bool
	TryLock(b []byte) bool
	ReleaseLock(b []byte) bool
	PurgeLock(int64)
	StoreContract(hash, value []byte, abi *ContractContext.ContractABI) bool
	GetContract(hash []byte) (string, bool)
	StoreFile(hash string, value []byte) bool
	RetrieveFile(hash string) []byte
	StoreKV(hash string, value []byte) bool
	RetrieveKV(hash string) []byte
	StoreTxToBlock(ts int64, hash string)
	StoreTx(hash string, txtype string, txTs int64)
	GetTxFromBlock(ts int64) []string
	GetTx(hash string) string
	GetContext(hash []byte) (map[string]string, *ContractContext.ContractABI, bool)
	CreateContext(hash []byte)
	StoreContext(hash []byte, ctx map[string]string)
}
