package cmd

import (
	"github.com/juliankoehn/wetter-service/cities"
	"github.com/juliankoehn/wetter-service/config"
	"github.com/spf13/cobra"
)

var refreshCitieseCmd = cobra.Command{
	Use:  "refresh-cities",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, refreshCities)
	},
}

func refreshCities(config *config.Configuration) {
	cities.Update(config)
}
