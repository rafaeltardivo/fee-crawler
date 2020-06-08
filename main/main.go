package main

import (
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rafaeltardivo/fee-crawler/api"
	"github.com/rafaeltardivo/fee-crawler/rates"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Clears cache on every server start
	rates.ClearCache()

	// Schedules invalidation for everyday at 14:30 UTC (or 16:30 CET)
	// This is the estimated time for daily update according to European Central Bank
	invalidateCacheTask := gocron.NewScheduler(time.UTC)
	invalidateCacheTask.Every(1).Day().At("14:30").Do(rates.ClearCache)
	invalidateCacheTask.StartAsync()
}

func main() {
	http.HandleFunc("/graphql", api.Serve)

	logger.Info("API server is ready to accept connections")
	http.ListenAndServe(":9000", nil)
}
