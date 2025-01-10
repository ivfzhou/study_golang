package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func handleHttp(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	ipPort := r.RemoteAddr
	ipPortArr := strings.Split(ipPort, ":")
	ip := ""
	if len(ipPortArr) != 0 {
		ip = ipPortArr[0]
	}
	reqBody, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(reqBody))
	if len(proxyAddr) != 0 {
		r.URL.Host = proxyAddr
		index := strings.Index(r.URL.Path, proxyPrefixPath)
		if index != -1 {
			r.URL.Path = r.URL.Path[len(proxyPrefixPath):]
		}
	}

	rsp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Println("ERROR", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer rsp.Body.Close()

	status := rsp.Status
	rspBody, _ := io.ReadAll(rsp.Body)
	if printAddr == ip && (strings.Contains(url, traceDomain) || traceDomain == "*") {
		toPrint(r.Method, url, string(reqBody), status, string(rspBody), r.Header, rsp.Header)
	}
	copyHeader(w.Header(), rsp.Header)
	w.WriteHeader(rsp.StatusCode)
	w.Write(rspBody)
}

var count int

func toPrint(method, url, reqBody, status, rspBody string, reqHeader, rspHeader map[string][]string) {
	count++
	fmt.Printf("[序号]：%d\n", count)
	fmt.Println(method, url, status)
	fmt.Print("[请求头]：")
	for k, vals := range reqHeader {
		fmt.Printf("%s=%s; ", k, strings.Join(vals, ","))
	}
	fmt.Println()
	fmt.Print("[响应头]：")
	for k, vals := range rspHeader {
		fmt.Printf("%s=%s; ", k, strings.Join(vals, ","))
	}
	fmt.Println()
	if len(reqBody) != 0 {
		fmt.Print("[请求体]：")
		fmt.Println(reqBody)
	}
	if len(rspBody) != 0 {
		fmt.Print("[响应体]：")
		fmt.Println(rspBody)
	}
	fmt.Println()
	fmt.Println()
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for _, v := range values {
			dst.Add(key, v)
		}
	}
}
