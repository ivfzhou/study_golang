package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

func main() {
	AnotherUsage()
}

func SampleUsage() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		w.Write([]byte("Hello World!"))
	})
	svr := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println(svr.ListenAndServeTLS("https/server.crt", "https/server.key"))
}

func AnotherUsage() {
	certificate, err := tls.LoadX509KeyPair("https/server.crt", "https/server.key")
	if err != nil {
		panic(err)
	}
	rootCAPool := x509.NewCertPool()
	caCert, err := os.ReadFile("https/ca.crt")
	if err != nil {
		panic(err)
	}
	rootCAPool.AppendCertsFromPEM(caCert)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		w.Write([]byte("Hello World!"))
	})
	svr := http.Server{
		Addr:    ":8080",
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    rootCAPool,
		},
	}
	fmt.Println(svr.ListenAndServeTLS("", ""))
}
