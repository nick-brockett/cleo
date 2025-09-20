package http

import (
	"net/http"

	"cleo.com/internal/core/domain/model"
	"cleo.com/internal/core/port"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewHealthMetricParserHandler(
	logger *logrus.Logger,
	parserService port.HealthMetricParserService,
) HealthMetricParserHandler {
	return HealthMetricParserHandler{
		logger:        logger,
		parserService: parserService,
	}
}

type HealthMetricParserHandler struct {
	logger        *logrus.Logger
	parserService port.HealthMetricParserService
}

func (h *HealthMetricParserHandler) Parse(c *gin.Context) {
	var note model.ClinicalNote

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	valid, err := note.Valid()
	if err != nil {
		h.logger.Infof("error encountered: invalid clinical note: %s", err.Error())
	}
	if !valid || err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	healthMetric, err := h.parserService.ParseClinicalNote(&note)
	if err != nil {
		h.logger.Infof("error encountered  for service to parse clinical note: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, healthMetric)

}
