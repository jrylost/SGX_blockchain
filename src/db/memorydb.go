package db

import (
	"SGX_blockchain/src/vm"
	"container/list"
	"encoding/hex"
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
	contextdb map[string]vm.StorageInterface
}

func InitMemorydb() *memorydb {
	d := &memorydb{}
	d.mdb = make(map[string][]byte)
	d.lock = make(map[string]*list.Element)
	d.lru = list.New()
	d.filedb = make(map[string][]byte)
	d.kvdb = make(map[string][]byte)
	d.contextdb = make(map[string]vm.StorageInterface)
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
	b := d.ReleaseLock(s)
	if !b {
		return false
	}
	d.mdb[key] = v
	return true
}

func (d *memorydb) StoreContract(hash, value []byte) bool {
	key := hex.EncodeToString(hash)
	d.filedb[key] = value
	return true
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
