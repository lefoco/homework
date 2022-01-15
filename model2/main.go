package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	_ "io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", rootHandler)
	serveMux.HandleFunc("/healthz", healthz)
	serveMux.HandleFunc("/debug/pprof/", pprof.Index)
	serveMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	serveMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	serveMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	httpServer := http.Server{
		Addr:    ":80",
		Handler: serveMux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Listening: %s\n", err)
		}
	}()

	fmt.Println("Http Server Listening on port:80...")
	<-done
	fmt.Println("Http Server exit...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		glog.Flush()
		cancel()
	}()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Println("Http Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Http Server exit...")

}

func rootHandler(response http.ResponseWriter, request *http.Request) {

	response.Write([]byte("hello httpserver..."))

	defer func() {
		err := recover()
		if err != nil {
			log.Fatal(err)
		} else {
			response.WriteHeader(http.StatusOK)
		}
	}()

	//接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range request.Header {
		for _, value := range v {
			response.Header().Set(k, value)
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	if os.Getenv("VERSION") == "" {
		os.Setenv("VERSION", "V1.0")
	}
	response.Header().Add("VERSION", os.Getenv("VERSION"))
	fmt.Println("VERSION =", os.Getenv("VERSION"))

	//Server 端记录访问日志包括客户端 IP，//HTTP 返回码, 输出到 server 端的标准输出
	fmt.Println("Client IP =", getClientIp())
	fmt.Println("http Status code =", http.StatusOK)
}

func healthz(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	fmt.Println("healthz return code:", http.StatusOK)
}

/**
Get Client IP
*/
func getClientIp() string {
	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range interfaceAddrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
