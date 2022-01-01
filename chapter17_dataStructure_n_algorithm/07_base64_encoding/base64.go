package main

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"unsafe"
)

var BASE64Table = "IJjkKLMNO567PQX12RVW3YZaDEFGbcdefghiABCHlSTUmnopqrxyz04stuvw89+/"

func Encode(data string) string {
	content := *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&data))))
	coder := base64.NewEncoding(BASE64Table)
	return coder.EncodeToString(content)
}

func Decode(data string) string {
	coder := base64.NewEncoding(BASE64Table)
	result, _ := coder.DecodeString(data)
	return *(*string)(unsafe.Pointer(&result))
}

func main() {

	strTest := "I love this beautiful world!"
	strEncrypt := Encode(strTest)
	strDecrypt := Decode(strEncrypt)
	fmt.Println("Encrypted:", strEncrypt)
	fmt.Println("Decrypted:", strDecrypt)

}
