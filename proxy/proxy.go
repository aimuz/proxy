package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

func HandlerInfo(writer http.ResponseWriter, modPath, version string) {
	info, err := executeGoCommandInfo(modPath, version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := os.Open(info.Info)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	_, err = io.Copy(writer, f)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

func HandlerZip(writer http.ResponseWriter, modPath, version string) {
	info, err := executeGoCommandInfo(modPath, version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := os.Open(info.Zip)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	_, err = io.Copy(writer, f)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

func HandlerMod(writer http.ResponseWriter, modPath, version string) {
	info, err := executeGoCommandInfo(modPath, version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := os.Open(info.GoMod)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()

	_, err = io.Copy(writer, f)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	return
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

	cmd.Dir = os.Getenv("GOPATH") + "/src"
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
