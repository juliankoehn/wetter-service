package storage

import (
	"net/url"

	"github.com/delivc/team/storage/namespace"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"    // default import
	_ "github.com/jinzhu/gorm/dialects/mysql"    // default import
	_ "github.com/jinzhu/gorm/dialects/postgres" // default import
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // default import
	"github.com/juliankoehn/wetter-service/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Connect opens a new Database Connection with Gorm
func Connect(config *config.Configuration) (*gorm.DB, error) {
	if config.DB.Driver == "" && config.DB.URL != "" {
		u, err := url.Parse(config.DB.URL)
		if err != nil {
			return nil, errors.Wrap(err, "parsing db connection url")
		}
		config.DB.Driver = u.Scheme

		if config.DB.Driver == "sqlite3" {
			config.DB.URL = u.Host
		}
	}

	db, err := gorm.Open(config.DB.Driver, config.DB.URL)
	if err != nil {
		return nil, errors.Wrap(err, "opening database connection")
	}
	if err := db.DB().Ping(); err != nil {
		return nil, errors.Wrap(err, "pinging database connection")
	}

	if config.DB.Namespace != "" {
		namespace.SetNamespace(config.DB.Namespace)
	}

	if logrus.StandardLogger().Level == logrus.DebugLevel {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	return db, nil
}
