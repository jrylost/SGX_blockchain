package dao

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

//打开数据库
func Open()  {
	var err error
	//数据存储路径和一些初始文件
	Db,err = leveldb.OpenFile("./levelDB/db",nil)
	if err != nil {
		log.Fatalln(err)
	}
}
