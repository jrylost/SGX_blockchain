package service

import (
	"bufio"
	"fmt"
	"levelDB_demo/src/dao"
	"log"
	"os"
	"strings"
)

func Recover(logName string){
	log.Println("prepare to Recover...")
	//按行读取日志
	inputLog, inputError := os.OpenFile(logName, os.O_RDWR, 0666)
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputLog\n")
		return
	}
	defer inputLog.Close()
	inputReader := bufio.NewReader(inputLog)
	for {
		//按行读取
		inputString, _ := inputReader.ReadString('\n')
		//查找执行成功的日志记录
		if(strings.Contains(inputString,"key")){
			//按空格分割
			arr:=strings.Fields(inputString)
			//寻找冒号的索引
			keyIndex:=strings.Index(arr[2],":")
			key:=arr[2][keyIndex+1:]
			valueIndex:=strings.Index(arr[3],":")
			value:=arr[3][valueIndex+1:]
			dao.Add(key,value)
		}
		if strings.Contains(inputString,"prepare to Recover..."){
			log.Println("Recover finished!")
			return
		}
	}
}
