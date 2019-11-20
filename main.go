package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aimuz/proxy/pkg/server"
	"github.com/aimuz/proxy/pkg/server/cmd"
	"github.com/spf13/cobra"
)

func main() {
	serverCmd := cmd.NewServerDefaultCommand()
	serverCmd.Run = run
	err := serverCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	rootCtx := SetupSignalHandler(context.Background())
	ctx := context.WithValue(rootCtx, "cmd", cmd)
	server.Start(ctx)
}

func SetupSignalHandler(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}
