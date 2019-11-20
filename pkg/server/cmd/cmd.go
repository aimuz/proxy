package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewServerDefaultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: filepath.Base(os.Args[0]) + " is goproxy server",
		Long:  `Private goproxy for Enterprise`,
	}

	flag := cmd.Flags()
	flag.String("address", ":8081", "goproxy server address")
	return cmd
}
