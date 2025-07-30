package router

import (
	"itinerary/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/generate", handler.GenerateItineraryPDF)

	return r
}
