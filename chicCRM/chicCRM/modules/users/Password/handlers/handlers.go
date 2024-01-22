package handlers

import (
	"chicCRM/modules/users/Password/models"
	"chicCRM/modules/users/Password/services"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type HandlerPort interface {
	InitPasswordChicCRMHandlers(c *gin.Context)
	ChangePasswordChicCRMHandlers(c *gin.Context)
	RequestResetPasswordHandlers(c *gin.Context)
	ResetPasswordChicCRMHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) InitPasswordChicCRMHandlers(c *gin.Context) {
	var initPassword models.InitPassword
	if err := c.ShouldBindJSON(&initPassword); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	email, ok := claims["username"].(string)
	if !ok || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Invalid token"})
		return
	}
	if email != initPassword.Username || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Invalid token or email"})
		return
	}
	existingToken := c.GetHeader("Authorization")
	if existingToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "No token provided"})
		return
	}
	existingToken = strings.TrimPrefix(existingToken, "Bearer ")
	if existingToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Invalid token format"})
		return
	}
	initPasswordAdditional := models.InitPasswordAdditional{
		ExistingToken: existingToken,
	}
	err := h.s.InitPasswordChicCRMServices(initPasswordAdditional, initPassword)
	if err != nil {
		switch err.Error() {
		case "password must not be empty and must be at least 8 characters long":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Password update successfully"})
}

func (h *handlerAdapter) ChangePasswordChicCRMHandlers(c *gin.Context) {
	var changePassword models.ChangePassword
	if err := c.ShouldBindJSON(&changePassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	err := h.s.ChangePasswordChicCRMServices(changePassword)
	if err != nil {
		switch err.Error() {
		case "password must not be empty and must be at least 8 characters long":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		case "incorrect old password":
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": err.Error()})
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Changed password successfully"})
}

func (h *handlerAdapter) RequestResetPasswordHandlers(c *gin.Context) {
	var requestResetPassword models.RequestResetPassword
	if err := c.ShouldBindJSON(&requestResetPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	fmt.Println(requestResetPassword)
	err := h.s.RequestResetPasswordChicCRMServices(requestResetPassword)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Email with the password reset link has been sent successfully"})
}

func (h *handlerAdapter) ResetPasswordChicCRMHandlers(c *gin.Context) {
	var resetPassword models.ResetPassword
	if err := c.ShouldBindJSON(&resetPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	claims := c.MustGet("claims").(jwt.MapClaims)
	email, ok := claims["username"].(string)
	if !ok || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Invalid token"})
		return
	}
	if email != resetPassword.Username || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Invalid token or email"})
		return
	}
	existingToken := c.GetHeader("Authorization")
	if existingToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "No token provided"})
		return
	}
	existingToken = strings.TrimPrefix(existingToken, "Bearer ")
	if existingToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": "Invalid token format"})
		return
	}
	ResetPasswordAdditional := models.InitPasswordAdditional{
		ExistingToken: existingToken,
	}
	err := h.s.ResetPasswordChicCRMServices(resetPassword, ResetPasswordAdditional)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		case "password must not be empty and must be at least 8 characters long":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Password reset successfully"})
}
