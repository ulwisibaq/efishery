package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/ulwisibaq/efishery/commodity/cmd/http"
)

var RootCmd = &cobra.Command{
	Use:   "efishery",
	Short: "efishery technical test",
}

func ExecuteCommodity() {

	RootCmd.AddCommand(http.HttpCmd)

	if err := RootCmd.Execute(); err != nil {
		log.Fatal("Error run auth service", err.Error())
		os.Exit(-1)
	}
}
