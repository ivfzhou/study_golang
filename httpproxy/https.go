package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func handleHttps(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)
	destConn, err := net.DialTimeout("tcp", r.Host, 60*time.Second)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		log.Println("ERROR", "not a Hijacker")
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
