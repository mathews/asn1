package asn1

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestBigIntMarshal(t *testing.T) {

	bIntString := "40144f0c70c094a2d4a48463655e45cf92596202886cb1b338b833cb8f601489"
	byteArray, err := hex.DecodeString(bIntString)
	if err != nil {
		fmt.Printf("error parse hex %s\n", err.Error())
		t.FailNow()
	}

	bInt := new(big.Int).SetBytes(byteArray)

	bbStr, err := Marshal(bInt)
	if err != nil {
		fmt.Printf("error Marshal %s\n", err.Error())
		t.FailNow()
	}
	var ret []byte
	if len(bbStr) > 32 && bbStr[2] == 0 {
		ret = bbStr[3:]
	} else {
		ret = bbStr[2:]
	}
	err = checkInteger(ret)
	if err != nil {
		fmt.Printf("asn1 Marshal %x as %x", bInt, ret)
		ret = bInt.Bytes()
	}
	if bytes.Compare(byteArray, ret) != 0 {
		fmt.Printf("not equal, bbStr %x, ret %x\n", bbStr, ret)
		t.FailNow()
	}
}
