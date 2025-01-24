package main

import (
	"math/rand"
	"time"
)

func main() {}

var randObj = rand.New(rand.NewSource(time.Now().UnixMilli() / 2))

func random(l int, chars string) string {
	bs := make([]byte, l)
	mask := len(chars)
	for i := range bs {
		bs[i] = chars[randObj.Int()%mask]
	}
	return string(bs)
}

const (
	characters      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charactersUpper = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers         = "0123456789"
)

func RandomChar(length int) string {
	return random(length, characters)
}

func RandomUpperChar(length int) string {
	return random(length, charactersUpper)
}

func RandomNumber(length int) string {
	return random(length, numbers)
}
