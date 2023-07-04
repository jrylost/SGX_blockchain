package server

import (
	"fmt"
	"reflect"
	"testing"
)

var abistring = `{
    "contract_name" : "evidence_v6",
    "contract_info" : "This is contract evidence_v6",
    "contract_functions" : [
        {
            "function_name" : "AddEvidence",
            "function_inputs" : [
                {
                    "input_name" : "context",
                    "input_type" : "context" 
                },
                {
                    "input_name" : "data",
                    "input_type" : "string" 
                }
            ],
            "function_outputs" : [
                {
                    "output_name" : "CreateResult",
                    "output_type" : "struct",
                    "output_components" : [
                        {
                            "component_name" : "evidenceId",
                            "component_type" : "string"
                        },
                        {
                            "component_name" : "txId",
                            "component_type" : "string"
                        },
                        {
                            "component_name" : "txTimestamp",
                            "component_type" : "string"
                        }
                    ] 
                }
            ]
        },
        {
            "function_name" : "QueryEvidenceById",
            "function_inputs" : [
                {
                    "input_name" : "context",
                    "input_type" : "context" 
                },
                {
                    "input_name" : "evidenceId",
                    "input_type" : "string" 
                }
            ],
            "function_outputs" : [
                {
                    "output_name" : "QueryResult",
                    "output_type" : "struct",
                    "output_components" : [
                        {
                            "component_name" : "key",
                            "component_type" : "string"
                        },
                        {
                            "component_name" : "record",
                            "component_type" : "string"
                        },
                        {
                            "component_name" : "bookmark",
                            "component_type" : "string"
                        }
                    ] 
                }
            ]
        },
        {
            "function_name" : "AddEvidence",
            "function_inputs" : [
                {
                    "input_name" : "context",
                    "input_type" : "context" 
                },
                {
                    "input_name" : "txId",
                    "input_type" : "string" 
                }
            ],
            "function_outputs" : [
                {
                    "output_name" : "QueryResult",
                    "output_type" : "struct",
                    "output_components" : [
                        {
                            "component_name" : "key",
                            "component_type" : "string"
                        },
                        {
                            "component_name" : "record",
                            "component_type" : "string"
                        },
                        {
                            "component_name" : "bookmark",
                            "component_type" : "string"
                        }
                    ] 
                }
            ]
        }
    ]
}`

var valuejsonstring = `[
    {
        "input_name" : "context",
        "input_type" : "context", 
        "input_value" : ""
    },
    {
        "input_name" : "data",
        "input_type" : "string",
        "input_value" : "{\"evidenceId\":\"aa\",\"uploaderSign\":\"cc\",\"content\":\"dd\"}"
    }
]`

func TestABIParser(t *testing.T) {

	res, err := ABIParser(abistring)
	if err != nil {
		t.Fatalf("error in parser")
	}
	if res.ContractName != "evidence_v6" {
		t.Fatalf("No contract name!")
	}
	if res.ContractInfo != "This is contract evidence_v6" {
		t.Fatalf("No contract info!")
	}
	if res.ContractInfo != "This is contract evidence_v6" {
		t.Fatalf("No contract info!")
	}
}

func TestContractInputHandler(t *testing.T) {
	inputs, err := ContractInputHandler(valuejsonstring)
	if err != nil {
		t.Fatalf("wrong inputs")
	}
	if inputs[0].InputName != "context" {
		t.Fatalf("wrong input name")
	}
}

func TestContractInputVerify(t *testing.T) {
	inputsJson := `[
    {
        "input_name" : "context",
        "input_type" : "context", 
        "input_value" : ""
    },
    {
        "input_name" : "data",
        "input_type" : "string",
        "input_value" : "{\"evidenceId\":\"aa\",\"uploaderSign\":\"cc\",\"content\":\"dd\"}"
    }
]`
	values, err := ContractInputHandler(inputsJson)
	if err != nil {
		t.Fatalf("wrong input")
	}
	res, err := ABIParser(abistring)
	if err != nil {
		t.Fatalf("wrong format")
	}
	var inputvalues []reflect.Value
	for _, function := range res.ContractFunctions {
		if function.FunctionName == "AddEvidence" {
			inputvalues, err = ContractInputVerify(values, function.FunctionInputs)
			if err != nil {
				t.Fatalf(err.Error())
			}
			break
		}
	}
	fmt.Println(inputvalues)
	fmt.Println(inputvalues[1])
	fmt.Println(inputvalues[1].Kind().String())
}
