package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	_ "io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", rootHandler)
	serveMux.HandleFunc("/healthz", healthz)

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

	defer func() {
		err := recover()
		if err != nil {
			log.Fatal(err)
		} else {
			response.WriteHeader(http.StatusOK)
		}
	}()

	fmt.Println("Client IP =", getIp())
	if os.Getenv("VERSION") != "" {
		response.Header().Set("VERSION", os.Getenv("VERSION"))
	}

	for k, v := range request.Header {
		for _, value := range v {
			response.Header().Set(k, value)
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	response.Write([]byte("hello httpserver..."))
}

func healthz(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	fmt.Println("healthz return code: ", http.StatusOK)
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
