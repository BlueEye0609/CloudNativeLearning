/*
1. 接收客户端 request，并将 request 中带的 header 写入 response header
2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4.当访问 localhost/healthz 时，应返回200
*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	for key, value := range req.Header {
		w.Header().Set(key, value[0])
	}

	// 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	var version string = os.Getenv("VERSION")
	w.Header().Add("version", version)

	// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	userAddr := req.RemoteAddr
	if strings.Contains(userAddr, ":") {
		fmt.Println("IP", net.ParseIP(strings.Split(userAddr, ":")[0]), "Response Code", http.StatusOK)
	} else {
		fmt.Println("IP", net.ParseIP(userAddr), "Response Code", http.StatusOK)
	}
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	// 4.当访问 localhost/healthz 时，应返回200
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.ListenAndServe(":80", nil)
}
