package vm

type vmContext struct {
	OwnerId          []byte //33字节
	Sender           string
	StorageInterface map[string]interface{}
}

type VirtualMachine struct {
	code map[string]string
}

func NewVirtualMachine() *VirtualMachine {
	code := make(map[string]string)
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

//func (v *VirtualMachine) Call(sender, contractName, funcName, hashValue string, params map[string][]byte) bool {
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
