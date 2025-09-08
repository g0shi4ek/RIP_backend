package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// creatorId получаем в хендлере и закидываем в слой репы?

func (h *ChargingHandler) DeleteChargingApplicationById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	// id, err := strconv.Atoi(c.Param("id"))
	//if err != nil{
	//	logrus.Error(err)
	//  	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": "Failed to get id",
	//	})
	//}
	applicationId := 1
	err := h.chargingRepository.DeleteChargingApplicationById(ctx, uint(applicationId))
	if err != nil {
		h.ErrorChargingHandler(c, http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusFound, "/tariffs")
}

func (h *ChargingHandler) GetChargingApplicationById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	//id, err := strconv.Atoi(c.Param("id"))
	//if err != nil {
	//    h.ErrorChargingHandler(c, http.StatusBadRequest, fmt.Errorf("invalid application ID"))
	//    return
	//}
	applicationId := uint(1)
	creatorId := uint(3) //  => из jwt

	fmt.Println(applicationId, creatorId)
	application, err := h.chargingRepository.GetChargingApplicationById(ctx, applicationId, creatorId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.ErrorChargingHandler(c, http.StatusNotFound, fmt.Errorf("application not found"))
		} else {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
		}
		return
	}

	chargingOrdersList, err := h.chargingRepository.GetChargingOrdersByApplicationId(ctx, applicationId)
	if err != nil {
		h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
		return
	}
	fmt.Println(application.Id)

	c.HTML(http.StatusOK, "application.html", gin.H{
		"application": application,
		"orders":      chargingOrdersList,
	})
}

func (h *ChargingHandler) AddTariffToApplication(c *gin.Context) {
	// добавить автосоздание заявки, если не существует

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	tariffId, err := strconv.Atoi(c.Param("tariffId"))
	if err != nil {
		h.ErrorChargingHandler(c, http.StatusBadRequest, fmt.Errorf("invalid tariff ID"))
		return
	}
	applicationId := uint(1)
	creatorId := uint(3) //  => из jwt

	application, err := h.chargingRepository.GetChargingApplicationById(ctx, applicationId, creatorId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.ErrorChargingHandler(c, http.StatusNotFound, fmt.Errorf("application not found"))
		} else {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, err)
		}
		return
	}
	fmt.Println(applicationId, tariffId)

	err = h.chargingRepository.CreateChargingOrder(c.Request.Context(), uint(tariffId), application.Id)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			h.ErrorChargingHandler(c, http.StatusBadRequest, fmt.Errorf("tariff already exists in application"))
		} else {
			h.ErrorChargingHandler(c, http.StatusInternalServerError, fmt.Errorf("failed to add tariff: %v", err))
		}
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/application/%d", applicationId))
}
