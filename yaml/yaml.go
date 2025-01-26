package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type T struct {
	Server *int   `yaml:"server"`
	Name   string `yaml:"name"`
}

func TestDecode() {
	var tt T
	file, err := os.ReadFile("./test.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, &tt)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tt.Server)

	tt = T{Name: "aa"}
	out, err := yaml.Marshal(tt)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", out)
}
