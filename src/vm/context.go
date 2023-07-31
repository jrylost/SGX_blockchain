package vm

import (
	"SGX_blockchain/src/vm/ContractContext"
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

func (v *VirtualMachine) Call(sender, txHashWith0x, contractName, funcName, code string, params []reflect.Value, ctx map[string]string) ([]interface{}, map[string]string, error) {
	var err error = nil

	i := interp.New(interp.Options{})
	context := ContractContext.Initial(ctx, txHashWith0x)
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
	evaluateFunction, err := i.Eval(contractName + "." + funcName)
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
	results := evaluateFunction.Call(params)

	var res []interface{}
	res = make([]interface{}, 0)
	for _, result := range results {
		res = append(res, result.Interface())
	}

	return res, context.Getstub().Strdb, err
}
