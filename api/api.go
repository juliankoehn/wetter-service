package api

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/jinzhu/gorm"
	"github.com/juliankoehn/wetter-service/config"
	"github.com/juliankoehn/wetter-service/storage"
	"github.com/juliankoehn/wetter-service/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// this is the actual api package
// we are launching our webServer here
// starting our updater and are serving data from openweathermap

// as we want to keep things as easy as possible we are using go-echo as webrouter
// it has a ton of built-in support for CORS / let's encrypt

// API holds our database, config and webServer
type API struct {
	e      *echo.Echo
	db     *gorm.DB
	config *config.Configuration
	sched  utils.ScheduledExecutor
	cache  *ristretto.Cache
}

// New creates a new API-Handler
func New(conf *config.Configuration) *API {
	a := &API{
		config: conf,
	}
	// open db
	db, err := storage.Connect(conf)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	// create a new echo instance
	e := echo.New()
	e.HideBanner = true

	// assign db, echo to API
	a.db = db
	a.e = e

	a.applyRoutes()
	if err := a.startCache(); err != nil {
		logrus.Fatalf("error starting cache: %+v", err)
	}

	a.sched = utils.NewScheduledExecutor(runnerDefaultTickTime, a.runner)
	// we should run the runner at least once before starting our
	// api services to "warm-up" our caches
	a.runner()

	return a
}

func (a *API) startCache() error {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: a.config.Cache.NumCounters,
		MaxCost:     a.config.Cache.MaxCosts,
		BufferItems: a.config.Cache.BufferItems,
		Metrics:     a.config.Cache.Metrics,
	})
	if err != nil {
		return err
	}

	a.cache = cache
	return nil
}

// Start our API Server
func (a *API) Start() {
	logrus.Info("Starting WebServer")
	var listenAddr string

	if a.config.Web.ListenAddr != "" {
		listenAddr = a.config.Web.ListenAddr
	} else {
		logrus.Info("missing ListenAddr in Config using default: `:1323`")
		listenAddr = ":1323"
	}

	// start our runner
	a.sched.Start()

	if a.config.Web.UseTLS {
		logrus.Info("TLS is enabled in Configuration, starting TLS-Server on Port 443")
		go func(c *echo.Echo) {
			if err := a.e.StartAutoTLS(":443"); err != nil {
				logrus.Info("shutting down the TLS server")
			}
			// a.e.Logger.Fatal(a.e.StartAutoTLS(":443"))
		}(a.e)
	}

	go func() {
		if err := a.e.Start(listenAddr); err != nil {
			logrus.Info("shutting down the http server")
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Stop(ctx); err != nil {
		logrus.Fatal(err)
	}
}

// Stop gracefull stop api server
func (a *API) Stop(ctx context.Context) error {
	a.cache.Close()
	a.sched.Stop()
	if err := a.db.Close(); err != nil {
		logrus.Fatal(err)
	}
	if err := a.e.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
