package db

import (
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
}

func InitMemorydb() *memorydb {
	d := &memorydb{}
	d.mdb = make(map[string][]byte)
	d.lock = make(map[string]*list.Element)
	d.lru = list.New()
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
