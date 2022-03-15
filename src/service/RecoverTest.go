package service

import (
	"fmt"
	"levelDB_demo/src/dao"
)

func RecoverTest(){
	//去掉0x
	value1:=dao.Get("6c467d5b0106ed5d352610d45644d4b5631a061b12acf501fcb38a869d84ba63")
	if value1!=nil {
		fmt.Printf("%s\n",value1)
	}else {
		fmt.Println("The value of key1 does not exist")
	}
	value2:=dao.Get("53459b42a452d320803cfa77b75e0d381fd7d604332dd662187bcd8d67dbb9d9")
	if value2!=nil {
		fmt.Printf("%s\n",value2)
	}else {
		fmt.Println("The value of key2 does not exist")
	}
}
