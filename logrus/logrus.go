package main

import (
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	SampleUse()
}

func SampleUse() {
	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              true,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              true,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			const moduleName = "golang"
			function = frame.Function
			file = frame.File
			if index := strings.LastIndex(frame.File, moduleName); index != -1 {
				file = frame.File[index+len(moduleName)+1:]
				function = frame.Function[strings.LastIndex(frame.Function, moduleName)+len(moduleName)+1:]
			}
			return
		},
	},
	/*&log.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			const moduleName = "golang"
			function = frame.Function
			file = frame.File
			if index := strings.LastIndex(frame.File, moduleName); index != -1 {
				file = frame.File[index+len(moduleName)+1:]
				function = frame.Function[strings.LastIndex(frame.Function, moduleName)+len(moduleName)+1:]
			}
			return
		},
		PrettyPrint: false,
	}*/)

	log.WithField("F1", "V1").Infof("hello")
	log.WithField("F1", "V1").Errorf("hello")
}
