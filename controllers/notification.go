package controllers

import (
	"net/http"
	"notifications/datarepo"
	"notifications/notifications"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
)

type NotificationRequest struct {
	UserID string `json:"userID"`
	Status bool   `json:"status"`
}

type FilledNotificationRequest struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}

func (hws *handlerWithServices) setNotificationStatus(g *gin.Context) {
	var (
		req NotificationRequest
		err error
	)

	if err = g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hws.Services.SetNotificationStatus(datarepo.ItemKey{ID: req.UserID, Kind: datarepo.KindDetails}, req.Status)

	if err != nil {
		g.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, nil)
}

func (hws *handlerWithServices) filledNotification(g *gin.Context) {
	var (
		req FilledNotificationRequest
		err error
	)

	if err = g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	details, err := hws.Services.GetUserByID(req.UserID)

	notifications.Send(details.DeviceToken, &notifications.Data{
		Alert: aws.String(req.Message),
		Sound: aws.String("default"),
		Badge: aws.Int(1),
	})

	if err != nil {
		g.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	g.JSON(http.StatusOK, nil)
}
