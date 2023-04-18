package accounts

import (
	"SGX_blockchain/src/db"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

var (
	accountid1 = "95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18"
	accountid2 = "95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b1f"
	filehash1  = "a4473b3f3a90025c936646d75195a3ab0a4685a31142423121375baed271dd6d"
)

func TestExternalAccount_Create(t *testing.T) {
	//初始化两个账户id
	id1, _ := hex.DecodeString(accountid1)
	id2, _ := hex.DecodeString(accountid2)

	//创建两个帐号
	ac1 := CreateAccount(id1, 100)
	ac2 := CreateAccount(id2, 0)

	//转账
	res, hashValue := ac1.Transfer(ac2, 50)

	//输出账户信息
	fmt.Println(res, hashValue)
	fmt.Println(ac1.Balance)
	fmt.Println(ac2.Balance)

	if ac1.Balance != 50 || ac2.Balance != 50 {
		t.Fatalf("After transfering,account of the two Account should be 50.0")
	}

}

func TestExternalAccount_Store(t *testing.T) {
	id := []byte{1, 2, 3}
	ac := CreateAccount(id, 100)
	res, _ := ac.MarshalMsg([]byte(""))

	newacc := NewAccount()
	newacc.UnmarshalMsg(res)

	jsonStr, _ := json.Marshal(newacc)
	fmt.Println(string(jsonStr))
}

func TestStoreFile(t *testing.T) {
	database := db.InitMemorydb()
	ac := NewAccount()
	id, _ := hex.DecodeString(accountid1)
	ac.Id = id
	//var lock bool
	var result bool
	var returnHash []byte
	if database.TryLock(ac.Id) {
		defer func() {
			database.ReleaseLock(ac.Id)
		}()
		fhash, _ := hex.DecodeString(filehash1)
		result, returnHash = ac.StoreFile(fhash)
	}

	res, _ := ac.MarshalMsg([]byte(""))

	newacc := NewAccount()
	newacc.UnmarshalMsg(res)
	jsonStr, _ := json.Marshal(newacc)
	fmt.Println(string(jsonStr))

	fmt.Println(hex.EncodeToString(returnHash), result)
}

func TestRecover(t *testing.T) {
	id := []byte{1, 2, 3}
	ac := CreateAccount(id, 100)
	res, _ := ac.MarshalMsg([]byte(""))
	ac2 := NewAccount()
	ac2.UnmarshalMsg(res)
	fmt.Println(len(res))
	fmt.Println(ac2.Id, ac2.File, ac2.Nonce, ac2.Balance)
}
