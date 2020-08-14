package cities

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/model"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/sirupsen/logrus"
)

var (
	tempFileName = "city.list.json.gz"
)

// package cities keeps the cities up2date from openweather

// Update updates the cities in our storage according to the given list of openweather
// as we currently dont know how we have to use the cities in future we are not deleting entries
func Update(conf *config.Configuration) {
	logrus.Info("Updating our cities")
	// the update endpoint is within our city list
	logrus.Infof("Downloading Update file from: %s", conf.OpenWeather.CityListEndpoint)
	if err := downloadUpdateFile(tempFileName, conf.OpenWeather.CityListEndpoint); err != nil {
		logrus.Fatal(err)
	}

	if ok := fileExists(tempFileName); !ok {
		logrus.Fatal("Downloaded file does not exists at target destination")
	}
	logrus.Info("Download Successfully")
	logrus.Info("Processing file")
	cities, err := processFile(tempFileName)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Cities have been successfully bound to struct")
	// currently we have a filtered list of cities by our country code
	// now we have to open a connection to our storage
	// and populate all data into it.
	db, err := storage.Connect(conf)
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()

	logrus.Info("Updating Cities in Storage")
	if err := model.UpdateOrCreateCities(db, cities); err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Cities have been Updated.")
}

func downloadUpdateFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// create the file on disk
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// determines if the file exists on our disk
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// processFile processes our source file
func processFile(srcFile string) ([]*model.City, error) {
	f, err := os.Open(srcFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(gzf)
	if err != nil {
		return nil, err
	}
	codeJSON := data

	cities := []*model.City{}

	if err := json.Unmarshal(codeJSON, &cities); err != nil {
		return nil, err
	}

	// return only german cities
	filtered := filterByCountryISO(cities, "DE")

	if err := os.Remove(srcFile); err != nil {
		logrus.Errorf("error deleting file %+v", err)
	}
	// skipping error, we dont care if the delete is successfully
	// it is just a 'cleanup' and not relevant to core functions
	return filtered, nil
}

// filterByCountryISO filters given struct map by iso code
func filterByCountryISO(cities []*model.City, iso string) []*model.City {
	filtered := []*model.City{}

	for _, city := range cities {
		if city.Country == iso {
			filtered = append(filtered, city)
		}
	}

	return filtered
}
