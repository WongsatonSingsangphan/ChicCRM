package handlers

import (
	"chicCRM/modules/queryData/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	GetCountriesHandlers(c *gin.Context)
	GetProvinceAmphoeTambonZipcodeHandlers(c *gin.Context)
	GetUniversalInfoHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) GetCountriesHandlers(c *gin.Context) {
	dataResponse, err := h.s.GetCountriesServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err})
		return
	}
	c.JSON(http.StatusOK, dataResponse)
}

func (h *handlerAdapter) GetProvinceAmphoeTambonZipcodeHandlers(c *gin.Context) {
	dataResponse, err := h.s.GetProvinceAmphoeTambonZipcodeServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err})
		return
	}
	c.JSON(http.StatusOK, dataResponse)
}

func (h *handlerAdapter) GetUniversalInfoHandlers(c *gin.Context) {
	dataResponse, err := h.s.GetUniversalInfoServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err})
		return
	}
	c.JSON(http.StatusOK, dataResponse)
}
