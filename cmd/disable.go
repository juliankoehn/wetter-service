package cmd

import (
	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/model"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// the enable command activates a city
// to be parsed by our updater
var disableCityCmd = cobra.Command{
	Use:  "disable",
	Long: "disable a city by id, use weather show to see a list of enabled cities",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, disableCity)
	},
}

func disableCity(config *config.Configuration) {
	if cityID == "" {
		logrus.Error("you must provide a city id: weather enable -i 123645")
	}
	// open db
	db, err := storage.Connect(config)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	city, err := model.FindEnabledCityByID(db, cityID)
	if err != nil {
		if model.IsNotFoundError(err) {
			logrus.Infof("City by id `%s` has not been found", cityID)
		} else {
			logrus.Fatal(err)
		}
		return
	}

	// delete it :O
	if err := db.Delete(city).Error; err != nil {
		logrus.Fatalf("Error deleting City: %+v", err)
	}

	logrus.Infof("Successfully deleted City: %s", city.ID)
}
