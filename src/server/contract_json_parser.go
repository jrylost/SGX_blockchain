package server

import (
	"errors"
	"fmt"
	"reflect"
)

//type FunctionInputsWithValue struct {
//	Inputs []FunctionInputWithValue
//}

type FunctionInputWithValue struct {
	InputName       string                     `json:"input_name"`
	InputType       string                     `json:"input_type"`
	InputValue      interface{}                `json:"input_value"`
	InputComponents []StructComponentWithValue `json:"input_components,omitempty"`
}

type StructComponentWithValue struct {
	ComponentName   string                     `json:"input_name"`
	ComponentType   string                     `json:"input_type"`
	InputValue      interface{}                `json:"input_value"`
	InputComponents []StructComponentWithValue `json:"input_components,omitempty"`
}

type StructComponent struct {
	ComponentName   string            `json:"component_name"`
	ComponentType   string            `json:"component_type"`
	InputComponents []StructComponent `json:"input_components,omitempty"`
}

type FunctionInput struct {
	InputName       string            `json:"input_name"`
	InputType       string            `json:"input_type"`
	InputComponents []StructComponent `json:"input_components,omitempty"`
}

type FunctionOutput struct {
	OutputName       string            `json:"output_name"`
	OutputType       string            `json:"output_type"`
	OutputComponents []StructComponent `json:"output_components,omitempty"`
}

type ContractFunction struct {
	FunctionName    string           `json:"function_name"`
	FunctionInputs  []FunctionInput  `json:"function_inputs,omitempty"`
	FunctionOutputs []FunctionOutput `json:"function_outputs,omitempty"`
}

type ContractABI struct {
	ContractName      string             `json:"contract_name"`
	ContractInfo      string             `json:"contract_info"`
	ContractFunctions []ContractFunction `json:"contract_functions"`
}

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ABIParser(ABIJsonString string) (*ContractABI, error) {
	abi := &ContractABI{}
	err := json.Unmarshal([]byte(ABIJsonString), &abi)
	if err != nil {
		return abi, errors.New("wrong abi format")
	}
	return abi, nil
}

func ContractInputHandler(inputJsonString string) ([]FunctionInputWithValue, error) {
	//inputs := &FunctionInputsWithValue{}
	var inputs []FunctionInputWithValue
	err := json.Unmarshal([]byte(inputJsonString), &inputs)
	if err != nil {
		fmt.Println(err)
		panic("input value error")
	}
	return inputs, err
}

func ContractInputVerify(values []FunctionInputWithValue, formats []FunctionInput) ([]reflect.Value, error) {
	res := make([]reflect.Value, 0)
	for i, format := range formats {
		cval := values[i]
		if cval.InputName != format.InputName {
			return nil, errors.New(cval.InputName + ":input_name mismatch")
		}
		if cval.InputType != format.InputType {
			return nil, errors.New(cval.InputName + ":input_type mismatch")
		}
		if format.InputName == "context" {
			//todo context替换
			res = append(res, reflect.ValueOf(nil))
			continue
		}
		if format.InputType == "struct" {
			component, _, err := StructInputVerify(cval.InputComponents, format.InputComponents)
			if err != nil {
				return []reflect.Value{component}, errors.New(cval.InputName + "." + err.Error())
			}
			res = append(res, component)
			continue
		}
		if reflect.TypeOf(cval.InputValue).Kind().String() != format.InputType {
			return nil, errors.New(cval.InputName + ":input_value type mismatch")
		}
		if reflect.TypeOf(cval.InputValue).Kind().String() == "int" {
			res = append(res, reflect.ValueOf(reflect.ValueOf(cval.InputValue).Int()))
		} else {
			res = append(res, reflect.ValueOf(cval.InputValue))
		}

	}
	return res, nil
}

func StructInputVerify(value []StructComponentWithValue, format []StructComponent) (reflect.Value, reflect.Type, error) {
	res := reflect.ValueOf(nil)
	structField := make([]reflect.StructField, 0)
	structValues := make([]reflect.Value, 0)
	names := make([]string, 0)

	for i, inputComponent := range format {
		valueComponent := value[i]
		if valueComponent.ComponentName != inputComponent.ComponentName {
			return res, nil, errors.New(valueComponent.ComponentName + ":input_name mismatch")
		}
		if valueComponent.ComponentType != inputComponent.ComponentType {
			return res, nil, errors.New(valueComponent.ComponentType + ":input_type mismatch")
		}
		switch valueComponent.ComponentType {
		case "context":
			return res, nil, errors.New(valueComponent.ComponentName + ":no context here")
		case "struct":
			component, structType, err := StructInputVerify(valueComponent.InputComponents, inputComponent.InputComponents)
			if err != nil {
				return res, nil, errors.New(valueComponent.ComponentName + "." + err.Error())
			} else {
				structField = append(structField, reflect.StructField{
					Name: inputComponent.ComponentName,
					Type: structType,
				})
				names = append(names, valueComponent.ComponentName)
				structValues = append(structValues, component)
			}
		case "int":
			if reflect.TypeOf(valueComponent.InputValue).Kind().String() != inputComponent.ComponentType {
				return res, nil, errors.New(valueComponent.ComponentType + ":input_value type mismatch")
			}
			names = append(names, valueComponent.ComponentName)
			structValues = append(structValues, reflect.ValueOf(reflect.ValueOf(valueComponent.InputValue).Int()))
		default:
			if reflect.TypeOf(valueComponent.InputValue).Kind().String() != inputComponent.ComponentType {
				return res, nil, errors.New(valueComponent.ComponentType + ":input_value type mismatch")
			}
			names = append(names, valueComponent.ComponentName)
			structValues = append(structValues, reflect.ValueOf(reflect.ValueOf(valueComponent.InputValue).Int()))
		}
	}

	returnType := reflect.StructOf(structField)
	returnValue := reflect.New(returnType).Elem()
	for i, name := range names {
		returnValue.FieldByName(name).Set(structValues[i])
	}
	return returnValue, returnType, nil
}
