package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
)

// cfb 密码反馈模式（Cipher FeedBack (CFB)）

// 加密。iv即初始向量，这里取密钥的前16位作为初始向量。
func encrypt(text string, key []byte) (string, error) {
	var iv = key[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, []byte(text))
	return hex.EncodeToString(encrypted), nil
}

// 解密
func decrypt(encrypted string, key []byte) (string, error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted), nil
}
func main() {
	str := "I love this beautiful world!"
	key := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}
	strEncrypted, err := encrypt(str, key)
	if err != nil {
		log.Fatal("Encrypted err:", err)
	}
	fmt.Println("Encrypted:", strEncrypted)
	strDecrypted, err := decrypt(strEncrypted, key)
	if err != nil {
		log.Fatal("Decrypted err:", err)
	}
	fmt.Println("Decrypted:", strDecrypted)
}

//Output:
//Encrypted: f866bfe2a36d5a43186a790b41dc2396234dd51241f8f2d4a08fa5dc
//Decrypted: I love this beautiful world!
