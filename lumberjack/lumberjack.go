package main

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {}

func SampleUse() {
	logrotater := lumberjack.Logger{
		Filename:   "lumberjack/test.log",
		MaxSize:    20,
		MaxAge:     90,
		MaxBackups: 50,
		LocalTime:  true,
		Compress:   false,
	}
	log.SetOutput(io.MultiWriter(&logrotater, os.Stdout))
	log.Println("logrotater.log", logrotater.Filename)
	err := logrotater.Rotate()
	if err != nil {
		log.Fatal(err)
	}
	err = logrotater.Close()
	if err != nil {
		panic(err)
	}
}
