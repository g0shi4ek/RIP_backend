package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
)

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

	if err != nil {
		h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
		return
	}

	applicationId := uint(1)
	creatorId := uint(3)

	application, _ := h.chargingRepository.GetChargingApplicationById(ctx, applicationId, creatorId)

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
		h.ErrorChargingHandler(c, http.StatusBadRequest, err)
		return
	}

	tariff, err := h.chargingRepository.GetTariffById(ctx, uint(id))
	if err != nil {
		h.ErrorChargingHandler(c, http.StatusBadRequest, err)
		return
	}

	c.HTML(http.StatusOK, "tariffDetails.html", gin.H{
		"tariff": tariff,
	})
}