package cmd

import (
	backend "github.com/anvh2/be-echo/services/echo"
	"github.com/spf13/cobra"
)

var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "Start echo service",
	Long:  `Start echo service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		server := backend.NewServer()
		return server.Run()
	},
}

func init() {
	RootCmd.AddCommand(echoCmd)
}
