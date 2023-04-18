package utils

import (
	"strings"
	"testing"
)

func TestGenerateRandomHexStringWith0x(t *testing.T) {
	s := GenerateRandomHexStringWith0x(32)
	if len(s) != 66 || !strings.HasPrefix(s, "0x") {
		t.Fatalf("format error")
	}
	//fmt.Println(s)
}
