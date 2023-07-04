package vm_test

import (
	"fmt"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
	"testing"
)

type Helloer interface {
	Hello()
}

func Hi(h Helloer) {
	println("In Hi:")
	h.Hello()
}

type Wrap struct {
	DoHello  func() // related to the Hello() method.
	Database map[string]interface{}
	// Other interface method wrappers...
}

func (w Wrap) Hello() { w.DoHello() }

func eval(t *testing.T, i *interp.Interpreter, src string) reflect.Value {
	t.Helper()
	res, err := i.Eval(src)
	if err != nil {
		t.Logf("Error: %v", err)
		if e, ok := err.(interp.Panic); ok {
			t.Logf(string(e.Stack))
		}
		t.FailNow()
	}
	return res
}

func TestInterface2(t *testing.T) {
	// export the Wrap type to the interpreter under virtual "wrap" package
	wraptest := &Wrap{}

	wraptest.Database = make(map[string]interface{})
	wraptest.Database["123"] = "456"

	i := interp.New(interp.Options{})
	err := i.Use(interp.Exports{
		"wrap/wrap": {
			"Wrap": reflect.ValueOf((*Wrap)(nil)),
			"s":    reflect.ValueOf(wraptest),
		},
	})

	err = i.Use(stdlib.Symbols)
	if err != nil {
		t.Fatal(err)
	}

	codesnippet :=
		`import "wrap"
		import "fmt"
	
    func useWrap (w *wrap.Wrap) string{
    	////w.Database["123"] = "6666"
    	//fmt.Println("???")
    	////println("???")
    	//test, _ := w.Database["123"]
    	////fmt.Println(test)
		test, _ := wrap.s.Database["123"].(string)
    	return test
}

	`
	i.Eval(codesnippet)
	//v, err := i.Eval("NewMyInt")

	//if err != nil {
	//	panic("wrong")
	//}

	//NewMyInt := v.Interface().(func(int) Wrap)
	//w := NewMyInt(4)
	//Hi(w)

	v2, err := i.Eval("useWrap")
	fmt.Println(err)
	useWrap := v2.Interface().(func(*Wrap) string)
	wraptest.Database["123"] = "777"
	//useWrap(wraptest)
	output := useWrap(wraptest)
	fmt.Println(output)

}

var depolyer = "95b01199edc2d8943ea9edb0ae5908a70bb960f23bc23310ed030e15ecc60b18"

//func TestInterface3(t *testing.T) {
//	contractName := "vote_example"
//	f, err := os.ReadFile("./examples/election/main.go")
//	if err != nil {
//		fmt.Println("read fail", err)
//	}
//	s := string(f)
//	params := make(map[string][]byte)
//	b := make([]byte, 8)
//	binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))
//	params["endTime"] = b
//
//	//fmt.Println(s)
//	virtualm := vm.NewVirtualMachine()
//	hashvalue := virtualm.Deploy(contractName, depolyer, s, params)
//
//	params = make(map[string][]byte)
//	params["proposal"] = []byte("1")
//	//fmt.Println(params)
//	virtualm.Call(depolyer, contractName, "Vote", hashvalue, params)
//	virtualm.Call(depolyer, contractName, "Vote", hashvalue, params)
//
//	virtualm.Call(depolyer, contractName, "Votes", hashvalue, params)
//
//}
