package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {}

func SampleUse() {
	claims := &jwt.RegisteredClaims{
		Issuer:    "ivfzhou",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := claim.SignedString([]byte("secret"))
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

	parse, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(parse.Raw)
	fmt.Println(parse.Header)
	fmt.Println(parse.Method)
	fmt.Println(parse.Valid)
	fmt.Println(parse.Claims)
}

func CreateAppleAPIJWT() {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &jwt.MapClaims{
		"iss": "issuer",
		"exp": time.Now().Unix() + int64(19*time.Minute.Seconds()),
		"aud": "appstoreconnect-v1",
		"iat": time.Now().Unix(),
	})
	token.Header["alg"] = "ES256"
	token.Header["kid"] = "keyId"
	token.Header["typ"] = "JWT"

	block, _ := pem.Decode([]byte(`secret.pkcs8`))
	if block == nil {
		panic("block is nil")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	signedString, err := token.SignedString(privateKey)
	if err != nil {
		panic(err)
	}

	fmt.Println(signedString)
}
