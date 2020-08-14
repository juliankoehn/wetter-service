package cmd

import (
	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/model"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var searchString = ""

var searchCmd = cobra.Command{
	Use:  "search",
	Long: "searches for a City",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, searchCity)
	},
}

func searchCity(config *config.Configuration) {
	if searchString == "" {
		logrus.Error("Missing `name` to search use: weather search -n cityName")
		return

	}
	db, err := storage.Connect(config)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	cities, err := model.FindCitiesByName(db, searchString)
	if err != nil {
		if !model.IsNotFoundError(err) {
			// if its not the 404 error we are returning a fatal
			// 404 is handled below.
			logrus.Fatal(err)
			return
		}
	}
	if len(cities) > 0 {
		logrus.Infof("Found the following cities by Criteria `%s`:", searchString)
		for _, city := range cities {
			logrus.Infof("ID: %s - Name: %s - Lat/Long: %f %f", city.ID, city.Name, city.Coords.Latitude, city.Coords.Longitude)
		}
	} else {
		logrus.Infof("could not find any cities by your search Criteria `%s`", searchString)
	}
}
