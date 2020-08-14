package cmd

import (
	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/model"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configFile = ""

var rootCmd = cobra.Command{
	Use: "weather",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

// RootCommand boots up our root command
// registers additional available commands to our cli-api
// and the config flag
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd, &refreshCitieseCmd, &searchCmd, &enableCityCmd, &disableCityCmd, &showCitiesCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "the config file to use")
	searchCmd.PersistentFlags().StringVarP(&searchString, "name", "n", "", "name of the city to search by")
	enableCityCmd.PersistentFlags().StringVarP(&cityID, "id", "i", "", "id of the city to be enabled")
	disableCityCmd.PersistentFlags().StringVarP(&cityID, "id", "i", "", "id of the city to be enabled")
	return &rootCmd
}

// execWithConfig
// loads the configuration from env or file
// and executes given command
func execWithConfig(cmd *cobra.Command, fn func(config *config.Configuration)) {
	config, err := config.LoadConfig(configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %+v", err)
	}
	if err := autoMigrate(config); err != nil {
		logrus.Fatal(err)
	}
	fn(config)
}

// autoMigrate uses the gorm automigrate function
// this func is chained in execWithConfig
func autoMigrate(config *config.Configuration) error {
	// open db
	db, err := storage.Connect(config)
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	if err := db.AutoMigrate(
		model.City{},
		model.CityEnabled{},
	).Error; err != nil {
		return err
	}

	return nil
}

func loadConfig() *config.Configuration {
	config, err := config.LoadConfig(configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %+v", err)
		return nil
	}

	return config
}
