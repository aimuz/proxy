package proxy

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/aimuz/go-json"
)

type Info struct {
	Path     string `json:"Path"`
	Version  string `json:"Version"`
	Info     string `json:"Info"`
	GoMod    string `json:"GoMod"`
	Zip      string `json:"Zip"`
	Dir      string `json:"Dir"`
	Sum      string `json:"Sum"`
	GoModSum string `json:"GoModSum"`
}

func Handler(fn func(info *Info) (io.ReadCloser, error)) func(writer http.ResponseWriter, modPath, version string) {
	return func(writer http.ResponseWriter, modPath, version string) {
		info, err := executeGoCommandInfo(modPath, version)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		r, err := fn(info)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Close()

		_, err = io.Copy(writer, r)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

type List struct {
	Version  string   `json:"Version"`
	Time     string   `json:"Time"`
	Versions []string `json:"Versions,omitempty"`
}

func HandlerList(writer http.ResponseWriter, modPath string) {
	list, err := executeGoCommandList(modPath)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(strings.Join(list.Versions, "\n")))
}

func executeGoCommand(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)

	fmt.Println(name, strings.Join(arg, " "))

	cmd.Dir = os.Getenv("GOPATH")
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return stdout, nil
}

// go mod download -json github.com/gliderlabs/logspout@v3.2.1+incompatible
func executeGoCommandInfo(modPath string, version string) (*Info, error) {
	b, err := executeGoCommand("go", "mod", "download", "-json", modPath+"@"+version)
	if err != nil {
		return nil, err
	}

	info := new(Info)
	return info, json.Unmarshal(b, info)
}

// go list -json -m -versions github.com/gliderlabs/logspout
func executeGoCommandList(modPath string) (*List, error) {
	b, err := executeGoCommand("go", "list", "-json", "-m", "-versions", modPath)
	if err != nil {
		return nil, err
	}

	list := new(List)
	return list, json.Unmarshal(b, list)
}
