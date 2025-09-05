package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/g0shi4ek/RIP_backend/config"
	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ChargingHandler struct {
	chargingRepository domain.IChargingRepository
	cfg                *config.Config
}

func NewChargingHandler(repo domain.IChargingRepository, cfg *config.Config) (*ChargingHandler, error) {
	return &ChargingHandler{
		chargingRepository: repo,
		cfg:                cfg,
	}, nil
}

func (h *ChargingHandler) InitRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")
	r.Static("/resources", "../resources")

	r.GET("/", h.GetAllTarrifs)
	r.GET("/tariff/:id", h.GetChargingTarrifById)
	r.GET("/application/:id", h.GetChargingApplicationById)
	r.GET("/search", h.SearchQuery)

	return r
}

func (h *ChargingHandler) GetAllTarrifs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	tariffsList, err := h.chargingRepository.GetAllTariffs(ctx)

	log.Println((*tariffsList)[0].Id)

	if err != nil || tariffsList == nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load tariffs",
		})
		return
	}
	application, _ := h.chargingRepository.GetChargingApplicationById(ctx, 1)
	log.Println((*application).Id)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"tariffs":     tariffsList,
		"application": application,
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

	c.HTML(http.StatusOK, "tariffDetails.tmpl", gin.H{
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

	c.HTML(http.StatusOK, "application.tmpl", gin.H{
		"application": application,
	})

}


func (h *ChargingHandler) SearchQuery(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
    defer cancel()
	
    searchQuery := c.Query("q")
    if searchQuery == "" {
        c.Redirect(http.StatusFound, "/")
        return
    }

    tariffsList, err := h.chargingRepository.SearchTariffs(ctx, searchQuery)
    if err != nil {
        logrus.Error(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Search failed",
        })
        return
    }

    application, _ := h.chargingRepository.GetChargingApplicationById(ctx, 1)
    
    c.HTML(http.StatusOK, "index.tmpl", gin.H{
        "tariffs":     tariffsList,
        "application": application,
        "searchQuery": searchQuery,
    })
}