package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	cloneHeaderToResponse(w, r)
	_, err := io.WriteString(w, os.Getenv("GOPROXY"))
	if err != nil {
		return
	}
	w.Header().Set("Go-Proxy", os.Getenv("GOPROXY"))
	log.Printf("clientIp: %v, statusCode: %v", getRequestIp(r), 200)
	fmt.Println(w.Header())
	w.WriteHeader(500)
}

func cloneHeaderToResponse(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		for _, headerValueItem := range v {
			w.Header().Set(k, headerValueItem)
		}
	}
}

func getRequestIp(r *http.Request) string {
	//这里，如果使用本地ip访问， r.RemoteAddr is [::1]:80, 无法分解，需要寻找合适的方法
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return ip
}

func main() {
	log.Printf("starting server...")
	http.HandleFunc("/healthz", healthz)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
