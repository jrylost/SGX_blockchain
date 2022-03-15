package service

import (
	"levelDB_demo/src/dao"
	"testing"
)

func BenchmarkPut(b *testing.B) {
	dao.Open()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Put("E:\\Go_code\\levelDB_demo\\data\\data.txt")
	}
	b.StopTimer()
}

