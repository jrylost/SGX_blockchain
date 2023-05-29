package crypto

import (
	"SGX_blockchain/src/utils"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"reflect"
)

func (k *KeyPair) SignJsonWithData(structBody interface{}) ([]byte, error) {
	if reflect.TypeOf(structBody).Kind() == reflect.Pointer {
		if reflect.ValueOf(structBody).Elem().Kind() != reflect.Struct {
			return []byte(""), errors.New("parameter structBody is not of struct type")
		}
	} else if reflect.TypeOf(structBody).Kind() != reflect.Struct {
		return []byte(""), errors.New("parameter structBody is not of struct type")
	}

	jsonBytes, err := json.Marshal(structBody)
	if err != nil {
		return []byte(""), err
	}

	data := gjson.GetBytes(jsonBytes, "data")
	if !data.Exists() {
		return []byte(""), errors.New("missing data field")
	}

	body := []byte(data.String())
	sig := k.SignMessage(body)

	jsonBytes, err = sjson.SetBytes(jsonBytes, "signature", utils.EncodeBytesToHexStringWith0x(sig))
	if err != nil {
		return []byte(""), err
	}
	return jsonBytes, nil

}
