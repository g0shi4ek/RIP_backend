package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ChargingHandler struct {
	chargingRepository domain.IChargingRepository
}

func NewChargingHandler(repo domain.IChargingRepository) (*ChargingHandler, error) {
	return &ChargingHandler{
		chargingRepository: repo,
	}, nil
}

func (h *ChargingHandler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.Static("/resources", "./resources")

	r.GET("/", h.GetTarrifs)
	r.GET("/tariff/:id", h.GetChargingTarrifById)
	r.GET("/application/:id", h.GetChargingApplicationById)
	
	return r
}

func (h *ChargingHandler) GetTarrifs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	var tariffsList *[]domain.ChargingTariff
	var err error

	searchQueryTariff := c.Query("tariffName")
	if searchQueryTariff == "" {
		tariffsList, err = h.chargingRepository.GetAllTariffs(ctx)
	} else {
		tariffsList, err = h.chargingRepository.SearchTariffs(ctx, searchQueryTariff)
	}

	if err != nil || tariffsList == nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load tariffs",
		})
		return
	}

	application, _ := h.chargingRepository.GetChargingApplicationById(ctx, 1)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"tariffs":      tariffsList,
		"application":  application,
		"searchTariff": searchQueryTariff,
	})
}

func (h *ChargingHandler) GetChargingTarrifById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get id",
		})
		return
	}
	tariff, err := h.chargingRepository.GetTariffById(ctx, id)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "There is no such tariff",
		})
		return
	}

	c.HTML(http.StatusOK, "tariffDetails.html", gin.H{
		"tariff": tariff,
	})
}

func (h *ChargingHandler) GetChargingApplicationById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	// id, err := strconv.Atoi(c.Param("id"))
	//if err != nil{
	//	logrus.Error(err)
	//  	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": "Failed to get id",
	//	})
	//}
	application, err := h.chargingRepository.GetChargingApplicationById(ctx, 1)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "There is no such application",
		})
		return
	}

	tariffsList, err := h.chargingRepository.GetAllTariffs(ctx)

	if err != nil || tariffsList == nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load tariffs",
		})
		return
	}
	orderTariffList := (*tariffsList)[4:]

	c.HTML(http.StatusOK, "application.html", gin.H{
		"application": application,
		"tariffs":     orderTariffList,
	})

}