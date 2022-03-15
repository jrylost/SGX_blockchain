package main

import (
	"levelDB_demo/src/dao"
	"levelDB_demo/src/service"
)

func main() {
	dao.Open()
	service.Put("./data/data.txt")
	//service.Put("E:\\Go_code\\levelDB_demo\\data\\data.txt")
	//service.PutTest()   // PutTest中包含了Put
	dao.Clean()
	service.Recover("./log.txt")
	// data_test
	service.RecoverTest()
}
