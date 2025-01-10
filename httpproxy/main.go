package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	listenAddr      string
	proxyAddr       string
	printAddr       string
	traceDomain     string
	proxyPrefixPath string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&listenAddr, "addr", ":8080", "监听地址；默认为本机8080端口")
	flag.StringVar(&proxyAddr, "proxyAddr", "", "本地代理请求地址")
	flag.StringVar(&printAddr, "printAddr", "127.0.0.1", "只打印该ip地址信息")
	flag.StringVar(&traceDomain, "traceDomain", "your_domain.com", "只打印该域名地址信息")
	flag.StringVar(&proxyPrefixPath, "proxyPrefixPath", "/path/prefix/", "本地代理请求地址前缀")
	flag.Parse()
}

func main() {
	log.Printf("listenAddr: %s; proxyAddr: %s; printAddr: %s; proxyPrefixPath: %s; traceDomain: %s\n",
		listenAddr, proxyAddr, printAddr, proxyPrefixPath, traceDomain)
	server := &http.Server{
		Addr: listenAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleHttps(w, r)
			} else {
				handleHttp(w, r)
			}
		}),
	}
	log.Println("WARN", server.ListenAndServe())
}
