package cmd

import (
	"github.com/juliankoehn/wetter-service/api"
	"github.com/juliankoehn/wetter-service/config"
	"github.com/spf13/cobra"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

func serve(config *config.Configuration) {
	a := api.New(config)
	a.Start()
}
