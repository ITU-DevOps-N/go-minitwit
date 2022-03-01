package controller

import (
	"github.com/gin-gonic/gin"
	// rename packe for easier read
	main "github.com/ITU-DevOps-N/go-minitwit"
	models "github.com/ITU-DevOps-N/go-minitwit/models"
)

// Get all Sample object from DB
func GetSamples(c *gin.Context) {
	// We access the DB pointer from the package
	var records []models.User

	var objects []models.Message

	results := main.DB.Find(&messages)

}
