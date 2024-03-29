package db

import (
	"SGX_blockchain/src/vm/ContractContext"
	"container/list"
	"encoding/hex"
	"encoding/json"
	"time"
)

const LockInterval = 30

type memorydb struct {
	mdb  map[string][]byte
	lock map[string]*list.Element
	lru  *list.List
	//Database
	filedb    map[string][]byte
	kvdb      map[string][]byte
	blockdb   map[int64][]string
	txdb      map[string]string
	abidb     map[string]*ContractContext.ContractABI
	contextdb map[string]map[string]string
}

func InitMemorydb() *memorydb {
	d := &memorydb{}
	d.mdb = make(map[string][]byte)
	d.lock = make(map[string]*list.Element)
	d.lru = list.New()
	d.filedb = make(map[string][]byte)
	d.kvdb = make(map[string][]byte)
	d.blockdb = make(map[int64][]string)
	d.txdb = make(map[string]string)
	d.abidb = make(map[string]*ContractContext.ContractABI)
	d.contextdb = make(map[string]map[string]string)
	return d
}

// PurgeLock Purge outdated lock
func (d *memorydb) PurgeLock(now int64) {
	for {
		ele := d.lru.Back()
		if ele != nil && now-ele.Value.(int64) > LockInterval {
			d.lru.Remove(ele)
		} else {
			break
		}
	}
}

// TryLock Try to acquire remote lock, return true if succeeded.
func (d *memorydb) TryLock(b []byte) bool {
	key := hex.EncodeToString(b)
	now := time.Now().Unix()
	d.PurgeLock(now)
	if _, ok := d.lock[key]; ok {
		return false
	} else {
		ele := d.lru.PushFront(now)
		d.lock[key] = ele
		return true
	}
	return false
}

// ReleaseLock Try to release remote lock, return false if remote lock purged.
func (d *memorydb) ReleaseLock(b []byte) bool {
	key := hex.EncodeToString(b)
	now := time.Now().Unix()
	d.PurgeLock(now)
	if val, ok := d.lock[key]; ok {
		d.lru.Remove(val)
		return true
	} else {
		return false
	}
	return false
}

func (d *memorydb) Get(b []byte) ([]byte, bool) {
	key := hex.EncodeToString(b)
	if val, ok := d.mdb[key]; ok {
		return val, ok
	}
	return []byte(""), false
}

func (d *memorydb) Put(s, v []byte) bool {
	key := hex.EncodeToString(s)
	//b := d.ReleaseLock(s)
	//if !b {
	//	return false
	//}
	d.mdb[key] = v
	return true
}

func (d *memorydb) StoreContract(hash, value []byte, abi *ContractContext.ContractABI) bool {
	key := hex.EncodeToString(hash)
	d.filedb[key] = value
	d.abidb[key] = abi
	return true
}

func (d *memorydb) GetContract(hash []byte) (string, bool) {
	key := hex.EncodeToString(hash)
	if val, exists := d.filedb[key]; exists {
		return string(val), true
	} else {
		return "", false
	}
}

func (d *memorydb) StoreFile(hash string, value []byte) bool {
	d.filedb[hash] = value
	return true
}

func (d *memorydb) RetrieveFile(hash string) []byte {
	if val, ok := d.filedb[hash]; ok {
		return val
	} else {
		return []byte("")
	}
}

func (d *memorydb) StoreKV(hash string, value []byte) bool {
	d.filedb[hash] = value
	return true
}

func (d *memorydb) RetrieveKV(hash string) []byte {
	if val, ok := d.filedb[hash]; ok {
		return val
	} else {
		return []byte("")
	}
}

func (d *memorydb) StoreTxToBlock(ts int64, hash string) {
	ts = ts / 1000
	if val, ok := d.blockdb[ts]; ok {
		d.blockdb[ts] = append(val, hash)
	} else {
		d.blockdb[ts] = []string{hash}
	}
}

func (d *memorydb) GetTxFromBlock(ts int64) []string {
	if val, ok := d.blockdb[ts]; ok {
		return val
	} else {
		return []string{}
	}
}

func (d *memorydb) StoreTx(hash string, txtype string, txTs int64) {
	resstruct := struct {
		Status string `json:"status"`
		Data   struct {
			Hash          string `json:"hash"`
			Type          string `json:"type"`
			TransactionTs int64  `json:"transactionTs"`
		} `json:"data"`
		Ts int64 `json:"ts"`
	}{
		Status: "ok",
		Data: struct {
			Hash          string `json:"hash"`
			Type          string `json:"type"`
			TransactionTs int64  `json:"transactionTs"`
		}{Hash: hash, Type: txtype, TransactionTs: txTs},
		Ts: 0,
	}
	resstr, _ := json.Marshal(resstruct)
	d.txdb[hash] = string(resstr)
}

func (d *memorydb) GetTx(hash string) string {
	if val, ok := d.txdb[hash]; ok {
		return val
	} else {
		resstruct := struct {
			Status string `json:"status"`
			Data   struct {
				Hash          string `json:"hash"`
				Type          string `json:"type"`
				TransactionTs int64  `json:"transactionTs"`
			} `json:"data"`
			Ts int64 `json:"ts"`
		}{
			Status: "error",
			Data: struct {
				Hash          string `json:"hash"`
				Type          string `json:"type"`
				TransactionTs int64  `json:"transactionTs"`
			}{Hash: hash, Type: "", TransactionTs: 0},
			Ts: 0,
		}
		resstr, _ := json.Marshal(resstruct)
		return string(resstr)
	}
}

func (d *memorydb) GetContext(hash []byte) (map[string]string, *ContractContext.ContractABI, bool) {
	key := hex.EncodeToString(hash)
	if val, ok := d.contextdb[key]; ok {
		if abival, ok2 := d.abidb[key]; ok2 {
			return val, abival, true
		}
	}
	return nil, nil, false
}

func (d *memorydb) CreateContext(hash []byte) {
	key := hex.EncodeToString(hash)
	d.contextdb[key] = make(map[string]string)
}

func (d *memorydb) StoreContext(hash []byte, ctx map[string]string) {
	key := hex.EncodeToString(hash)
	d.contextdb[key] = ctx
}
