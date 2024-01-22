package handlers

import (
	"chicCRM/modules/migrateData/models"
	"chicCRM/modules/migrateData/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	MigrateMemberDataByOrganizeHandlers(c *gin.Context)
	MigrateTeamleadTracByOrganizeHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) MigrateMemberDataByOrganizeHandlers(c *gin.Context) {
	var organize models.Organize
	if err := c.ShouldBindJSON(&organize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	err := h.s.MigrateMemberDataByOrganizeServices(organize)
	if err != nil {
		switch err.Error() {
		case "please insert valid organize ID":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		case "organize_id not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Migrated Successfully"})
}

func (h *handlerAdapter) MigrateTeamleadTracByOrganizeHandlers(c *gin.Context) {
	var organize models.Organize
	if err := c.ShouldBindJSON(&organize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	fmt.Println(organize)
	err := h.s.MigrateTeamleadTracByOrganizeServies(organize)
	if err != nil {
		switch err.Error() {
		case "please insert valid organize ID":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		case "organize_id not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Migrated Successfully"})
}
