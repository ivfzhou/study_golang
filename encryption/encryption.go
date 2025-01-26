package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"log"
)

func main() {}

func SHA256() {
	sha := sha256.New()
	if _, err := sha.Write([]byte("12345678")); err != nil {
		log.Fatal(err)
	}
	res := string(sha.Sum(nil))
	fmt.Printf("%X", res)
}

func SHA1() {
	sha := sha1.New()
	if _, err := sha.Write([]byte("12345678")); err != nil {
		log.Fatal(err)
	}
	res := string(sha.Sum(nil))
	fmt.Printf("%X", res)
}

func MD5() {
	sha := md5.New()
	if _, err := sha.Write([]byte("12345678")); err != nil {
		log.Fatal(err)
	}
	res := string(sha.Sum(nil))
	fmt.Printf("%X", res)
}

func AES() {
	var (
		key  = []byte("woyaozouqi16weia")
		text = []byte("加密文本16字节整数倍。")
		iv   = []byte("yeyaozouqi16weia")
	)

	if len(key)%aes.BlockSize != 0 {
		panic("密钥位数不合法")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	if block.BlockSize() != len(iv) {
		panic("密钥和iv长度不同")
	}
	encrypter := cipher.NewCBCEncrypter(block, iv)
	crypto := make([]byte, len(text))
	encrypter.CryptBlocks(crypto, text)
	fmt.Printf("%X\n", crypto)

	decrypter := cipher.NewCBCDecrypter(block, iv)
	res := make([]byte, len(text))
	decrypter.CryptBlocks(res, crypto)
	fmt.Printf("%s\n", res)
}
