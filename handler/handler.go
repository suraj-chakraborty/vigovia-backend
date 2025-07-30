package handler

import (
	"itinerary/model"
	service "itinerary/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateItineraryPDF(c *gin.Context) {
	var itinerary model.ItineraryData
	if err := c.ShouldBindJSON(&itinerary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filePath, err := service.GeneratePDF(itinerary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Failed to generate PDF": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pdfPath": filePath})
}
