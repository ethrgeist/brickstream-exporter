package main

import (
	"github.com/ethrgeist/brickstream-exporter/controller"
	"github.com/ethrgeist/brickstream-exporter/internal/database"
	"github.com/ethrgeist/brickstream-exporter/internal/logger"
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/ethrgeist/brickstream-exporter/repository"
	"github.com/ethrgeist/brickstream-exporter/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

var log zerolog.Logger

func main() {
	db := database.DB
	log.Info().Msg("Starting Brickstream Exporter")

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	siteRepository := repository.NewSiteRepository(db, log)
	brickstreamService := service.NewBrickstreamService(siteRepository, log)

	controller.NewBrickstreamController(router, brickstreamService, log)

	router.POST("/ingest/metrics", func(c *gin.Context) {
		var m models.MetricsV5

		if err := c.ShouldBindXML(&m); err != nil {
			c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Info().Str("Device", m.DeviceID).Msg("Received metrics")
		m.Process()
		log.Info().
			Str("Timezone", m.Properties.TimezoneParsed.String()).
			Time("DateLocal", m.ReportData.Reports[0].DateLocal).
			Time("DateUTC", m.ReportData.Reports[0].DateUTC).
			Time("StartTimeLocal", m.ReportData.Reports[0].Objects[0].Counts[0].StartTimeLocal).
			Time("StartTimeUTC", m.ReportData.Reports[0].Objects[0].Counts[0].StartTimeUTC).
			Time("EndTimeLocal", m.ReportData.Reports[0].Objects[0].Counts[0].EndTimeLocal).
			Time("EndTimeUTC", m.ReportData.Reports[0].Objects[0].Counts[0].EndTimeUTC).
			Time("UnixStartTimeParsed", m.ReportData.Reports[0].Objects[0].Counts[0].UnixStartTimeParsed).
			Int("Enters", m.ReportData.Reports[0].Objects[0].Counts[0].Enters).
			Int("Exits", m.ReportData.Reports[0].Objects[0].Counts[0].Exits).
			Msg("Processed metrics")
		c.Status(http.StatusOK)
	})

	err := router.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func init() {
	log = logger.GetLogger()
	err := database.DbConn()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
}
