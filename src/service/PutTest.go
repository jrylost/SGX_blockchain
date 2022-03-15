package service

import (
	"fmt"
	"time"
)

func PutTest() {
	t1:=time.Now()
	Put("E:\\Go_code\\levelDB_demo\\data\\data.txt")
	t2:=time.Now()
	fmt.Print("将data.txt存入数据库所消耗的时间为：")
	fmt.Println(t2.Sub(t1))
}
