package vm

import (
	"SGX_blockchain/src/vm/ContractContext"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
)

type VirtualMachine struct {
	code string
}

func NewVirtualMachine(code string) *VirtualMachine {
	return &VirtualMachine{
		code: code,
	}

}

// Deploy is used to deploy a contract
//func (v *VirtualMachine) Deploy(sender, contractName, funcName, hashValue string, params map[string][]byte) error {
//	i := interp.New(interp.Options{})
//	context := ContractContext.Initial()
//	stub := context.Getstub()
//	err := i.Use(interp.Exports{
//		"ContractContext/ContractContext": {
//			"Context": reflect.ValueOf(context),
//			"Stub":    reflect.ValueOf(stub),
//		},
//	})
//
//	err = i.Use(stdlib.Symbols)
//	if err != nil {
//		panic("wrong!")
//	}
//	src, _ := v.code[hashValue]
//	_, err = i.Eval(src)
//	//fmt.Println(src, err)
//	funcv, err := i.Eval(contractName + "." + funcName)
//	//fmt.Println(err, contractName+"."+funcName)
//	//fmt.Println(i, si)
//	callFune := funcv.Interface().(func(map[string][]byte) bool)
//	res := callFune(params)
//	return true
//}

func (v *VirtualMachine) Call(sender, contractName, funcName, code string, params []reflect.Value, ctx map[string]string) ([]interface{}, map[string]string, error) {
	var err error = nil

	i := interp.New(interp.Options{})
	context := ContractContext.Initial(ctx)
	stub := context.Getstub()
	err = i.Use(interp.Exports{
		"ContractContext/ContractContext": {
			"Context": reflect.ValueOf(*context),
			"Stub":    reflect.ValueOf(*stub),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	err = i.Use(stdlib.Symbols)
	if err != nil {
		return nil, nil, err
	}
	_, err = i.Eval(code)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Println(src, err)
	funcv, err := i.Eval(contractName + "." + funcName)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Println(err, contractName+"."+funcName)
	//fmt.Println(i, si)
	for i2, param := range params {
		if !param.IsValid() {
			params[i2] = reflect.ValueOf(context)
		}
	}
	results := funcv.Call(params)

	var res []interface{}
	res = make([]interface{}, 0)
	for _, result := range results {
		res = append(res, result.Interface())
	}
	//fmt.Println(stub.GetStringState("aa"))
	//fmt.Println(*context)
	fmt.Println(context.Getstub().Strdb)
	fmt.Println(context.Getstub().Strdb["Evi_aa"], "?????")
	return res, context.Getstub().Strdb, err
}
