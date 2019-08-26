package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aimuz/proxy/proxy"
)

type Handler func(writer http.ResponseWriter, mod, version string)

var addr = ":8081"

func main() {
	http.Handle("/", proxy.NewProxy())
	server := &http.Server{Addr: addr, Handler: http.DefaultServeMux}
	go func() {
		log.Println("Listen", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	initSignal(func(signal os.Signal) {
		log.Println("signal:", signal)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	})
}

func initSignal(cancel func(signal os.Signal)) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	c := <-quit
	cancel(c)
}
