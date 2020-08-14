package cmd

import (
	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/model"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cityID string

// the enable command activates a city
// to be parsed by our updater
var enableCityCmd = cobra.Command{
	Use:  "enable",
	Long: "enables a city by id, use weather search -n cityName to find it's ID",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, enableCity)
	},
}

func enableCity(config *config.Configuration) {
	if cityID == "" {
		logrus.Error("you must provide a city id: weather enable -i 123645")
	}
	// open db
	db, err := storage.Connect(config)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	city, err := model.FindCityByID(db, cityID)
	if err != nil {
		if model.IsNotFoundError(err) {
			logrus.Infof("City by id `%s` has not been found", cityID)
		} else {
			logrus.Fatal(err)
		}
		return
	}

	// just copy oneonone city -> cityEnabled
	enabled := model.CityEnabled(*city)
	if err := model.EnableCity(db, enabled); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("City `%s` with ID: `%s` has been enabled.", enabled.Name, enabled.ID)
}
