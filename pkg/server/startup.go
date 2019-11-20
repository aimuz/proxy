package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aimuz/proxy/pkg/proxy"
	"github.com/spf13/cobra"
)

func Start(ctx context.Context) {
	cmd, ok := ctx.Value("cmd").(*cobra.Command)
	if !ok {
		log.Panic("server start fail")
	}

	flags := cmd.Flags()
	address, err := flags.GetString("address")
	if err != nil {
		log.Panic(fmt.Errorf("get address flag: %w", err))
	}

	http.Handle("/", proxy.NewProxy())
	server := &http.Server{Addr: address, Handler: http.DefaultServeMux}
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	log.Println("Listen", address)
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
