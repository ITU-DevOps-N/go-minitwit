package controller

import (
  "net/http"
  "github.com/gin-gonic/gin"
  // rename packe for easier read  
  models "github.com/ITU-DevOps-N/go-minitwit/models"
  main "github.com/ITU-DevOps-N/go-minitwit"
)
// Get all Sample object from DB
func GetSamples(c *gin.Context) {
  // We access the DB pointer from the package 
  var records []Group

  var objects []models.Message

  results :=main.DB.Find(&messages)
  

}