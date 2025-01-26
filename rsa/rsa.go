package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log"
)

func main() {
	var (
		prk = `-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJTuVwgkvMR45TFM
/wOF7cRYQvSb3ugresf7qbjsfEUoBpOKTjhj5HWHxbaewAuT/2SO+3y4f/NMFyyM
ebFs4AILYsgvkn9lzRoAqPCA9qTo2WUixWcBP60f/xXqPpHb1waMRy2Gy+i4OOIR
6m4DzGGtYhpKRiaqwExsd1IKO2mxAgMBAAECgYEAj3H9C8/urUJQZqrlmOwfdhUY
8GdNTMvMJ+CCuaW1kBqcMvFso627N2S9j0babIxw2ddJ7Pf77UflrjfjYnweSQGv
6AjWguqN6JlV0WcSjfcbTHvn6Z9XUGmdUwi+Z7t3tXw2FzohBW81TG8yBjziTg40
7bHjfpl+al3a8AUB9h0CQQDD6wbzDBIRBU+pWyF296c3oyt9oScpfY/e8B9XzDJN
OQ5jIxeU4pt7qUaCjDl4zdSVUT7J7kVynoqdV3laKvszAkEAwpp98HPWWqMsB9yW
F1Ov/BjPz1vHvGdW4ELV3AmVWtRkLVxIM2RBnc1z2pz97D+1T6GfHVVy0H3aAvtW
st7niwJAYTgafbcKrAmPq0F+jLN99fzxUukKLuuQ3hcX5pB8kZdzjTxXslj0wNuS
EqwUxN6W0/W6C6hCLAuCS2uh212iwwJBAKSyLY3X630QBc6tgJVDbXh040NCENvB
tcPcrLQppC6X/CRrqmtcGTBdVgSZw0kzbdZ0GX6w95e+O0k0v95oShsCQEPXViUk
iSk2kKZENxMXogBfdELcH9UJKsrjWenEiEAWkHNVo3/DvYJCB8aYB5F1bYq8EhPQ
695QNaJSLOSTY4E=
-----END PRIVATE KEY-----`
		puk = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCU7lcIJLzEeOUxTP8Dhe3EWEL0
m97oK3rH+6m47HxFKAaTik44Y+R1h8W2nsALk/9kjvt8uH/zTBcsjHmxbOACC2LI
L5J/Zc0aAKjwgPak6NllIsVnAT+tH/8V6j6R29cGjEcthsvouDjiEepuA8xhrWIa
SkYmqsBMbHdSCjtpsQIDAQAB
-----END PUBLIC KEY-----`
	)

	prk, puk, err := GenerateKeyPair(1024)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(prk)
	log.Println(puk)
	str, err := EncodeRsaStr("hello world", puk)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(str)
	str, err = DecodeRSAStr(str, prk)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(str)
}

func GenerateKeyPair(bits int) (prk string, puk string, err error) {
	prkObj, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	prkBytes, err := x509.MarshalPKCS8PrivateKey(prkObj)
	if err != nil {
		return "", "", err
	}
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: prkBytes,
	}
	prk = string(pem.EncodeToMemory(&block))

	pukBytes, err := x509.MarshalPKIXPublicKey(&prkObj.PublicKey)
	if err != nil {
		return "", "", err
	}
	block = pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pukBytes,
	}
	puk = string(pem.EncodeToMemory(&block))
	return
}

func DecodeRSAStr(text, prk string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(prk))
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pkey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("parse prk error, not a *rsa.PrivateKey")
	}
	res, err := rsa.DecryptPKCS1v15(rand.Reader, pkey, bs)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func EncodeRsaStr(text, pukStr string) (string, error) {
	block, _ := pem.Decode([]byte(pukStr))
	pukInter, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	puk, ok := pukInter.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("parse puk error, not a *rsa.PublicKey")
	}
	res, err := rsa.EncryptPKCS1v15(rand.Reader, puk, []byte(text))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}
