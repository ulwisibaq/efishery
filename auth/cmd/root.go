package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/ulwisibaq/efishery/auth/cmd/http"
)

var RootCmd = &cobra.Command{
	Use:   "efishery",
	Short: "efishery technical test",
}

func ExecuteAuth() {

	RootCmd.AddCommand(http.HttpCmd)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Error run auth service", err.Error())
		os.Exit(-1)
	}
}
