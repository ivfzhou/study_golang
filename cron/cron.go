package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {}

type logger struct{}

func (l *logger) Info(msg string, keysAndValues ...any) {
	fmt.Println(msg, keysAndValues)
}

func (l *logger) Error(err error, msg string, keysAndValues ...any) {
	fmt.Println(msg, keysAndValues)
}

func SampleUse() {
	c := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local), cron.WithLogger(&logger{}))
	_, _ = c.AddFunc("* * * * * *", func() {
		fmt.Println("once second")
	})
	_, _ = c.AddFunc("@hourly", func() {
		fmt.Println("once hour")
	})
	_, _ = c.AddFunc("@every 1h30m", func() {
		fmt.Println("once hour and thirty minutes")
	})
	c.Run()
	// c.Start()
	// c.Stop()
}
