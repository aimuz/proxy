package proxy

import (
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/mod/module"
)

func NewProxy() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_paths := strings.Split(request.URL.Path, "/@v/")
		if len(_paths) != 2 {
			http.Redirect(writer, request, "https://github.com/aimuz/proxy", http.StatusFound)
			return
		}

		version := _paths[1]
		_path, err := module.UnescapePath(_paths[0][1:])
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		var handler func(writer http.ResponseWriter, modPath, version string)
		switch {
		case strings.HasSuffix(version, ".info"):
			version = strings.TrimSuffix(version, ".info")
			handler = Handler(func(info *Info) (reader io.ReadCloser, e error) {
				return os.Open(info.Info)
			})
		case strings.HasSuffix(version, ".mod"):
			version = strings.TrimSuffix(version, ".mod")
			handler = Handler(func(info *Info) (reader io.ReadCloser, e error) {
				return os.Open(info.GoMod)
			})
		case strings.HasSuffix(version, ".zip"):
			version = strings.TrimSuffix(version, ".zip")
			handler = Handler(func(info *Info) (reader io.ReadCloser, e error) {
				return os.Open(info.Zip)
			})
		case strings.EqualFold(version, "list"):
			HandlerList(writer, _path)
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
}
