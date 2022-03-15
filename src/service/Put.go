package service

import (
	"bufio"
	"fmt"
	"io"
	"levelDB_demo/src/dao"
	"log"
	"os"
	"strings"
)

//将文件名为fileName的文件存入数据库
func Put(fileName string) {
	//0777表示：创建了一个普通文件，所有人拥有所有的读、写、执行权限
	//
	//0666表示：创建了一个普通文件，所有人拥有对该文件的读、写权限，但是都不可执行
	//
	//0644表示：创建了一个普通文件，文件所有者对该文件有读写权限，用户组和其他人只有读权限，都没有执行权限
	inputFile, inputError := os.OpenFile(fileName, os.O_RDWR, 0777)
	if inputError != nil {
		fmt.Println(inputError)
		return
	}
	defer inputFile.Close()
	inputReader := bufio.NewReader(inputFile)
	for {
		//按行读取
		inputString, readerError := inputReader.ReadString('\n')
		//按空格分割成key和value，并去掉开头的0x
		arr:=strings.Fields(inputString)
		//最后一行的len为0
		if len(arr)!=0 {
			key:=arr[0][2:]
			value:=arr[1][2:]
			dao.Add(key,value)
		}
		if readerError == io.EOF {
			log.Println("Put finished!")
			return
		}
	}
}
