package main

import (
	"net/http"
	"strings"

	"github.com/aimuz/proxy/proxy"
	"golang.org/x/mod/module"
)

type Handler func(writer http.ResponseWriter, mod, version string)

func main() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		_paths := strings.Split(request.URL.Path, "/@v/")
		if len(_paths) != 2 {

		}

		modPath := _paths[0][1:]
		version := _paths[1]

		_path, err := module.UnescapePath(modPath)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		handler := func(writer http.ResponseWriter, mod, version string) {
			return
		}

		switch {
		case strings.HasSuffix(version, ".info"):
			version = strings.TrimSuffix(version, ".info")
			handler = proxy.HandlerInfo
		case strings.HasSuffix(version, ".mod"):
			version = strings.TrimSuffix(version, ".mod")
			handler = proxy.HandlerMod
		case strings.HasSuffix(version, ".zip"):
			version = strings.TrimSuffix(version, ".zip")
			handler = proxy.HandlerZip
		case strings.EqualFold(version, "list"):
			proxy.HandlerList(writer, modPath)
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
