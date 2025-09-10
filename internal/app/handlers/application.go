package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// creatorId получаем в хендлере и закидываем в слой репы?

func (h *ChargingHandler) DeleteChargingApplicationById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId := uint(3)
	application, err := h.chargingRepository.GetChargingApplicationByStatus(ctx, creatorId, "draft")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.ErrorChargingHandler(c, http.StatusNotFound, fmt.Errorf("application not found"))
		} else {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
		}
		return
	}
	err = h.chargingRepository.DeleteChargingApplicationById(ctx, uint(application.Id))
	if err != nil {
		h.ErrorChargingHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusFound, "/tariffs")
}

func (h *ChargingHandler) GetChargingApplicationByStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId := uint(3) //  => из jwt
	application, err := h.chargingRepository.GetChargingApplicationByStatus(ctx, creatorId, "draft")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.ErrorChargingHandler(c, http.StatusNotFound, fmt.Errorf("application not found"))
		} else {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
		}
		return
	}

	chargingOrdersList, err := h.chargingRepository.GetChargingOrdersByApplicationId(ctx, application.Id)
	if err != nil {
		h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "application.html", gin.H{
		"application": application,
		"orders":      chargingOrdersList,
	})
}

func (h *ChargingHandler) AddTariffToApplication(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	tariffId, err := strconv.Atoi(c.Param("tariffId"))
	if err != nil {
		h.ErrorChargingHandler(c, http.StatusBadRequest, fmt.Errorf("invalid tariff ID"))
		return
	}
	creatorId := uint(3) //  => из jwt

	application, err := h.chargingRepository.GetChargingApplicationByStatus(ctx, creatorId, "draft")
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
			return
		}
		application, err = h.chargingRepository.CreateDraftChargingApplication(ctx, creatorId)
		if err != nil {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
			return
		}
	}

	err = h.chargingRepository.CreateChargingOrder(ctx, uint(tariffId), application.Id)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			h.ErrorChargingHandler(c, http.StatusBadRequest, fmt.Errorf("tariff already exists in application"))
		} else {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, fmt.Errorf("failed to add tariff: %v", err))
		}
		return
	}

	c.Redirect(http.StatusFound, "/tariffs")
}
