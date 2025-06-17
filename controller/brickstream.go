package controller

import (
	"github.com/ethrgeist/brickstream-exporter/models"
	"github.com/ethrgeist/brickstream-exporter/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

type BrickstreamController interface {
}

type brickstreamController struct {
	brickstreamService service.BrickstreamService
	log                zerolog.Logger
}

func NewBrickstreamController(engine *gin.Engine, brickstreamService service.BrickstreamService, log zerolog.Logger) {
	controller := &brickstreamController{
		brickstreamService: brickstreamService,
		log:                log,
	}

	engine.POST("/api/v1/brickstream/ingest/xml", controller.ParseMetricsXML)

	log.Debug().Msg("Brickstream Controller initialized")

}

func (bc *brickstreamController) ParseMetricsXML(c *gin.Context) {
	var m models.MetricsV5

	if err := c.ShouldBindXML(&m); err != nil {
		c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().Str("Device", m.DeviceID).Msg("Received metrics")
	m.Process()
	bc.log.Info().
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

	err := bc.brickstreamService.SaveMetrics(&m)
	if err != nil {
		bc.log.Error().Err(err).Msg("Failed to save metrics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metrics"})
		return
	}

	c.Status(http.StatusOK)
}
