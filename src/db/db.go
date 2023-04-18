package db

type Database interface {
	Get(s []byte) ([]byte, bool)
	Put(s, v []byte) bool
	TryLock(b []byte) bool
	ReleaseLock(b []byte) bool
	PurgeLock(int64)
}
