package vote_example

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)
import "vmcontext"

func Create(params map[string][]byte) bool {
	t, _ := params["endTime"]
	endTime := int64(binary.LittleEndian.Uint64(t))
	vmcontext.StorageInterface["time"] = endTime
	vmcontext.StorageInterface["1"] = 0
	vmcontext.StorageInterface["2"] = 0
	vmcontext.StorageInterface["3"] = 0
	return true
}

func Vote(params map[string][]byte) bool {
	endTime, _ := vmcontext.StorageInterface["time"]
	e := endTime.(int64)
	p, _ := params["proposal"]

	proposal, _ := strconv.Atoi(string(p))
	if time.Now().Unix() > e {
		return false
	}

	if proposal > 3 || proposal < 1 {
		return false
	}
	if i, _ := vmcontext.StorageInterface[vmcontext.Sender]; i == true {
		return false
	}

	currentVotes, _ := vmcontext.StorageInterface[string(p)]
	tv := currentVotes.(int) + 1
	vmcontext.StorageInterface[string(p)] = tv

	return true
}

func Votes(params map[string][]byte) bool {
	currentVotes1, _ := vmcontext.StorageInterface["1"].(int)
	currentVotes2, _ := vmcontext.StorageInterface["2"].(int)
	currentVotes3, _ := vmcontext.StorageInterface["3"].(int)
	fmt.Println("总票数：", currentVotes1+currentVotes2+currentVotes3)
	fmt.Println("1号票数：", currentVotes1)
	fmt.Println("2号票数：", currentVotes2)
	fmt.Println("3号票数：", currentVotes3)
	return true
}
