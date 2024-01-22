package handlers

import (
	"chicCRM/modules/users/login/models"
	"chicCRM/modules/users/login/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	LoginChicCRMHandlers(c *gin.Context)
	LoginTeamleadSecuredocHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) LoginChicCRMHandlers(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	tokenJWT, err := h.s.LoginChicCRMServices(loginRequest)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		case "username or password is invalid":
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": err.Error()})
		case "username must be a valid email address":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Login successfully", "token": tokenJWT})
}

func (h *handlerAdapter) LoginTeamleadSecuredocHandlers(c *gin.Context) {
	var loginRequestTeamlead models.LoginRequestTeamlead
	if err := c.ShouldBindJSON(&loginRequestTeamlead); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	tokenJWT, err := h.s.LoginTeamleadSecuredocServices(loginRequestTeamlead)
	if err != nil {
		switch err.Error() {
		case "username must be a valid email address":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		case "username or password does not match":
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Login successfully", "token": tokenJWT})
}
