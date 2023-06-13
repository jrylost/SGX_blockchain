package vm

import (
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
