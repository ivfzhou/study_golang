package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {}

type User struct {
	Name        string        `json:"name"`
	Has         bool          `json:"has,string"`           // string：在序列化与反序列化时，都要是字符串格式
	Age         int           `json:"age,omitempty,string"` // omitempty：为零值时，序列化时不处理
	Ignore      string        `json:"-"`                    // -：忽略序列化与反序列化
	Birthday    time.Time     `json:"birthday"`             // 在序列化与反序列化时，处理 time.RFC3339 字符串格式
	WaitingTime time.Duration `json:"waitingTime"`          // 在序列化与反序列化时，处理数字格式。单位纳秒
}

func Marshal() {
	var js = &User{
		Name:        "ivzhou",
		Has:         false,
		Age:         18,
		Birthday:    time.Now(),
		WaitingTime: time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond,
		Ignore:      "ohoo",
	}
	bs, err := json.MarshalIndent(js, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}

func Unmarshal() {
	bs := []byte(`{
    "name": "ivzhou",
    "has": "false",
    "age": "18",
    "birthday": "2025-01-09T14:57:07.6169215+08:00",
    "waitingTime": 3661001001001
}`)
	var js User
	err := json.Unmarshal(bs, &js)
	if err != nil {
		panic(err)
	}
	fmt.Println(js)
}

func UnMarshalTime() {
	bs := []byte(`{
    "birthday": "2025-01-09T14:57:07+08:00"
}`)
	var js User
	err := json.Unmarshal(bs, &js)
	if err != nil {
		panic(err)
	}
	fmt.Println(js)

	t, err := time.Parse(time.RFC3339, "2025-01-09T14:57:07+08:00")
	if err != nil {
		panic(err)
	}
	fmt.Println(t)
}
