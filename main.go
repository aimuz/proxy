package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aimuz/proxy/proxy"
	"golang.org/x/mod/module"
)

type Handler func(writer http.ResponseWriter, mod, version string)

func main() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		_paths := strings.Split(request.URL.Path, "/@v/")
		if len(_paths) != 2 {
			http.NotFound(writer, request)
			return
		}

		modPath := _paths[0][1:]
		version := _paths[1]

		_path, err := module.UnescapePath(modPath)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		var handler Handler
		switch {
		case strings.HasSuffix(version, ".info"):
			version = strings.TrimSuffix(version, ".info")
			handler = proxy.Handler(func(info *proxy.Info) (reader io.ReadCloser, e error) {
				f, err := os.Open(info.Info)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					return
				}
				return ioutil.NopCloser(f), err
			})
		case strings.HasSuffix(version, ".mod"):
			version = strings.TrimSuffix(version, ".mod")
			handler = proxy.Handler(func(info *proxy.Info) (reader io.ReadCloser, e error) {
				f, err := os.Open(info.GoMod)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					return
				}
				return ioutil.NopCloser(f), err
			})
		case strings.HasSuffix(version, ".zip"):
			version = strings.TrimSuffix(version, ".zip")
			handler = proxy.Handler(func(info *proxy.Info) (reader io.ReadCloser, e error) {
				f, err := os.Open(info.Zip)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusBadRequest)
					return
				}
				return ioutil.NopCloser(f), err
			})
		case strings.EqualFold(version, "list"):
			proxy.HandlerList(writer, modPath)
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

	http.ListenAndServe(":8081", http.DefaultServeMux)
}
