package vm

import (
	"SGX_blockchain/src/crypto"
	"SGX_blockchain/src/utils"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
)

type vmContext struct {
	OwnerId          []byte //33字节
	Sender           string
	StorageInterface map[string]interface{}
}

//go:generate msgp

type StorageInterface struct {
	ByteStorage    map[string][]byte  `msg:"byte"`
	Int64Storage   map[string][]int64 `msg:"int"`
	StringStorage  map[string]string  `msg:"string"`
	Float64Storage map[string]float64 `msg:"float"`
}

type VirtualMachine struct {
	code    map[string]string
	storage map[string]StorageInterface
}

func NewVirtualMachine() *VirtualMachine {
	code := make(map[string]string)
	storage := make(map[string]StorageInterface)

	return &VirtualMachine{
		code:    code,
		storage: storage,
	}

}

func (v *VirtualMachine) Deploy(contractName, deployer string, src string, params map[string][]byte) string {
	i := interp.New(interp.Options{})
	si := StorageInterface{}
	err := i.Use(interp.Exports{
		"vmcontext/vmcontext": {
			"StorageInterface": reflect.ValueOf(si),
			"Sender":           reflect.ValueOf(deployer),
		},
		//"standardlib/standardlib": {
		//	"keccak256": reflect.ValueOf(crypto.Keccak256),
		//},
	})

	err = i.Use(stdlib.Symbols)
	if err != nil {
		panic("wrong!")
	}

	_, err = i.Eval(src)
	fmt.Println(err)
	funcv, err := i.Eval(contractName + ".Create")
	fmt.Println(err)
	fmt.Println(err, contractName+".Create")
	createfunc := funcv.Interface().(func(map[string][]byte) bool)
	res := createfunc(params)
	if res != true {
		panic("wrong")
	}
	hashBytes := crypto.Keccak256([]byte(src))
	hashValue := utils.EncodeBytesToHexStringWith0x(hashBytes)
	v.storage[hashValue] = si
	v.code[hashValue] = src
	return hashValue
	//useWrap(wraptest)
	//output := useWrap(wraptest)
	//fmt.Println(hashValue, res, v)

}

func (v *VirtualMachine) Call(sender, contractName, funcName, hashValue string, params map[string][]byte) bool {
	i := interp.New(interp.Options{})
	si, _ := v.storage[hashValue]
	err := i.Use(interp.Exports{
		"vmcontext/vmcontext": {
			"StorageInterface": reflect.ValueOf(si),
			"Sender":           reflect.ValueOf(sender),
		},
	})

	err = i.Use(stdlib.Symbols)
	if err != nil {
		panic("wrong!")
	}
	src, _ := v.code[hashValue]
	_, err = i.Eval(src)
	//fmt.Println(src, err)
	funcv, err := i.Eval(contractName + "." + funcName)
	//fmt.Println(err, contractName+"."+funcName)
	//fmt.Println(i, si)
	callFune := funcv.Interface().(func(map[string][]byte) bool)
	res := callFune(params)
	v.storage[hashValue] = si
	//useWrap(wraptest)
	//output := useWrap(wraptest)
	fmt.Println(res)
	return true
}
