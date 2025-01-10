package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {}

func Connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     "ivfzhoudebian:6379",
		Username: "",
		Password: "123456",
		DB:       0,
	})
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	if err = client.Close(); err != nil {
		panic(err)
	}
}

func ConnectURL() {
	url, err := redis.ParseURL("redis://:123456@ivfzhoudebian:6379/0")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(url)
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	if err = client.Close(); err != nil {
		panic(err)
	}
}

func TXPipeline() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "123456",
	})
	ctx := context.Background()
	rsp, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)

	txPipeline := rdb.TxPipeline()
	txPipeline.HExists(ctx, "hash", "field")
	txPipeline.HGet(ctx, "hash0", "field0")
	res, err := txPipeline.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		panic(err)
	}
	fmt.Println(res)
}
