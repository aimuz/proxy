package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aimuz/proxy/proxy"
	"golang.org/x/mod/module"
)

type Handler func(writer http.ResponseWriter, mod, version string)

var addr = ":8081"

func main() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_paths := strings.Split(request.URL.Path, "/@v/")
		if len(_paths) != 2 {
			http.Redirect(writer, request,"https://github.com/aimuz/proxy",http.StatusFound)
			return
		}

		version := _paths[1]
		_path, err := module.UnescapePath(_paths[0][1:])
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		var handler Handler
		switch {
		case strings.HasSuffix(version, ".info"):
			version = strings.TrimSuffix(version, ".info")
			handler = proxy.Handler(func(info *proxy.Info) (reader io.ReadCloser, e error) {
				return os.Open(info.Info)
			})
		case strings.HasSuffix(version, ".mod"):
			version = strings.TrimSuffix(version, ".mod")
			handler = proxy.Handler(func(info *proxy.Info) (reader io.ReadCloser, e error) {
				return os.Open(info.GoMod)
			})
		case strings.HasSuffix(version, ".zip"):
			version = strings.TrimSuffix(version, ".zip")
			handler = proxy.Handler(func(info *proxy.Info) (reader io.ReadCloser, e error) {
				return os.Open(info.Zip)
			})
		case strings.EqualFold(version, "list"):
			proxy.HandlerList(writer, _path)
			return
		default:
			http.NotFound(writer, request)
			return
		}

		_version, err := module.UnescapeVersion(version)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		handler(writer, _path, _version)
	})

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
			log.Println("signal:", signal)
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
