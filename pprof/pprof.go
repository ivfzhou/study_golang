package main

import (
	"net/http"
	"net/http/pprof"
)

func main() {
	SampleUse()
}

func SampleUse() {
	mux := http.NewServeMux()
	mux.HandleFunc("/x/pprof", pprof.Profile)
	http.ListenAndServe(":8080", mux)
}
