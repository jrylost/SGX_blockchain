package dao

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"testing"
)

//打开数据库
func init()  {
	var err error
	//数据存储路径和一些初始文件
	Db,err = leveldb.OpenFile("./levelDB/db",nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestAdd(t *testing.T) {
	Add("1","a")
	got,_ := Db.Get([]byte("1"),nil)
	if string(got)!="a"{
		t.Errorf("got:%s,but want:a",got)
	}
}

func TestRemove(t *testing.T) {
	Add("1","a")
	Remove("1")
	got,_ := Db.Get([]byte("1"),nil)
	if string(got)!=""{
		t.Errorf("got:%s,but want:nil",got)
	}
}

func TestUpdate(t *testing.T) {
	Add("1","a")
	Update("1","b")
	got,_ := Db.Get([]byte("1"),nil)
	if string(got)!="b"{
		t.Errorf("got:%s,but want:b",got)
	}
}

func TestGet(t *testing.T) {
	Add("1","a")
	got:=Get("1")
	if string(got)!="a"{
		t.Errorf("got:%s,but want:a",got)
	}
}
