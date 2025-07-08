package controller

import (
	"github.com/ethrgeist/brickstream-exporter/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsController interface {
	MetricsHandler() gin.HandlerFunc
}

type metricsController struct {
	metricsService service.MetricsService
	log            zerolog.Logger
}

func NewMetricsController(engine *gin.Engine, ms service.MetricsService, log zerolog.Logger) MetricsController {
	controller := &metricsController{
		metricsService: ms,
		log:            log,
	}

	engine.GET("/metrics", controller.MetricsHandler())

	controller.log.Debug().Msg("Metrics Controller initialized")

	return controller
}

func (mc *metricsController) MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		mc.metricsService.UpdateSiteMetrics()
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	}
}
