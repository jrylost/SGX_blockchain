package dao

import (
	"fmt"
	"github.com/imroc/biu"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
)
func init() {
	//以 可读可写|没有时创建|文件尾部追加 的形式打开这个文件
	//os.O_RDONLY - 仅供读取使用
	//os.O_WRONLY - 仅供写入
	//os.O_RDWR - 开放阅读和写入
	logFile, err := os.OpenFile("./log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
}


var Db *leveldb.DB

func Add(key string,value string)  {
	log.Println("prepare to Add...")
	Db.Put([]byte(biu.BytesToBinaryString([]byte(key))),[]byte(biu.BytesToBinaryString([]byte(value))),nil)
	log.Printf("Add-->key:%s   value:%s\n",key,value)
}

func Remove(key string){
	log.Println("prepare to Remove...")
	value,_:=Db.Get([]byte(biu.BytesToBinaryString([]byte(key))),nil)
	Db.Delete([]byte(biu.BytesToBinaryString([]byte(key))), nil)
	log.Printf("Remove-->key:%s   value:%s\n",key,string(value))
}

func Update(key string,value string){
	log.Println("prepare to Update...")
	Db.Put([]byte(biu.BytesToBinaryString([]byte(key))),[]byte(biu.BytesToBinaryString([]byte(value))),nil)
	log.Printf("Update-->key:%s   value:%s\n",key,value)
}

func Get(key string) []byte  {
	log.Println("prepare to Get...")
	value,err := Db.Get([]byte(biu.BytesToBinaryString([]byte(key))),nil)
	if err != nil {
		return nil
	}
	log.Printf("Get-->key:%s   value:%s\n",[]byte(biu.BytesToBinaryString([]byte(key))),value)
	return value
}

func GetAll(){
	iter := Db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("%s,%s",key,value)
	}
	iter.Release()
}

func Clean(){
	log.Println("prepare to Clean...")
	iter := Db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		Db.Delete([]byte(biu.BytesToBinaryString([]byte(key))), nil)
	}
	iter.Release()
	log.Println("Clean finished!")
}