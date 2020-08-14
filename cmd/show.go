package cmd

import (
	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/model"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// show cmd displays all enabled cities
var showCitiesCmd = cobra.Command{
	Use:  "show",
	Long: "show returns a list of enabled cities",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, showCities)
	},
}

func showCities(config *config.Configuration) {
	// open db
	db, err := storage.Connect(config)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	cities, err := model.FindEnabledCities(db)
	if err != nil {
		logrus.Fatal(err)
	}

	if len(cities) > 0 {
		logrus.Info("The following Cities are enabled:")
		for _, city := range cities {
			logrus.Infof("ID: %s - Name: %s - Lat/Long: %f %f", city.ID, city.Name, city.Coords.Latitude, city.Coords.Longitude)
		}
	} else {
		logrus.Info("no cities enabled.")
	}

	_ = cities
}
