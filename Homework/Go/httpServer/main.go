package main

import (
	"context"
	"fmt"
	"httpServerDemo/metrics"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func index(w http.ResponseWriter, req *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	// 0 ms to 2s delay
	delay := rand.Intn(2000)
	time.Sleep(time.Millisecond * time.Duration(delay))

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
	// subscribe SIGINT/SIGTERM signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	metrics.Register()

	mux := http.NewServeMux()
	// 06.debug
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		// server connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("start http server failed, error: %s\n", err.Error())
		}
	}()

	<-stopChan // wait for SIGINT or SIGTERM
	log.Println("Shutting down server")

	// shut down gracefully, but wait no longer 30 seconds before halting.
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(ctx)
}
