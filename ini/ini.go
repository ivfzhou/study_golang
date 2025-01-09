package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ini/ini"
)

func main() {}

type Server struct {
	Name string `ini:"name"`
	IP   string `ini:"ip"`
}

func Load() {
	cfg, err := ini.Load([]byte(`; comment
name=ivfzhou ; oh comment
ip=127.0.0.1
`))
	if err != nil {
		panic(err)
	}
	section := cfg.Section(ini.DEFAULT_SECTION)
	fmt.Println(section.Comment)
	key := section.Key("name")
	fmt.Println(key.Value())
	fmt.Println(key.Comment)
	fmt.Println(key.Validate(func(s string) string {
		fmt.Println(s)
		return s + "_modified"
	}))
}

func MapToStruct() {
	cfg := Server{}
	err := ini.MapTo(&cfg, []byte(`name=ivfzhou
ip=127.0.0.1
`))
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}

func StructToFile() {
	cfg := Server{
		Name: "ivfzhou",
		IP:   "127.0.0.1",
	}
	file := ini.Empty()
	err := file.ReflectFrom(&cfg)
	if err != nil {
		panic(err)
	}
	n, err := file.WriteTo(os.Stdout)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
}

func ToFile() {
	file := ini.Empty()
	section := file.Section(ini.DEFAULT_SECTION)
	section.Comment = "comment"
	section.Key("name").Comment = "comments"
	section.Key("name").SetValue("ivfzhou")
	section.Key("ip").SetValue("127.0.0.1")
	n, err := file.WriteTo(os.Stdout)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
}

func ParseDuration() {
	var cfg struct {
		Duration time.Duration `ini:"duration"`
	}
	err := ini.MapTo(&cfg, []byte(`duration=24h`))
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.Duration)
}

func ParseBool() {
	var cfg struct {
		Bool bool `ini:"bool"`
	}
	err := ini.MapTo(&cfg, []byte(`bool=true`))
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.Bool)
}

func ParseTime() {
	var cfg struct {
		Time time.Time `ini:"time"`
	}
	err := ini.MapTo(&cfg, []byte(`time=2025-01-04T16:25:00+08:00"`))
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.Time)
	fmt.Println(time.Parse(time.RFC3339, "2025-01-04T16:25:00+08:00"))
}
