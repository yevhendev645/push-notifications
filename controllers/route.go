package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlerWithServices struct {
	Services Services
}

func NewHTTPHandler() http.Handler {
	hws := handlerWithServices{Services: &defaultServices{}}
	router := gin.Default()

	api := router.Group("/api/v1")
	api.POST("/set-notification-status", hws.setNotificationStatus)
	api.POST("/filled-notification", hws.filledNotification)

	return router
}
