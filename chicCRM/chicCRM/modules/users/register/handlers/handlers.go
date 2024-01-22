package handlers

import (
	"chicCRM/modules/users/register/models"
	"chicCRM/modules/users/register/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	RegisterChicCRMHandlers(c *gin.Context)
	ValidateDomainChicCRMHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) RegisterChicCRMHandlers(c *gin.Context) {
	var loginData models.RegisterRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	companyID, err := h.s.RegisterChicCRMServices(loginData)
	if err != nil {
		switch err.Error() {
		case "email already exists", "mobile Phone already exists", "mobile phone must be 10 digits 089-XXX-XXXX", "username must be a valid email address":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "companyID": companyID.CompanyID})
}

func (h *handlerAdapter) ValidateDomainChicCRMHandlers(c *gin.Context) {
	var validateDomainRequest models.ValidateDomainRequest
	if err := c.ShouldBindJSON(&validateDomainRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	// fmt.Println(validateDomainRequest)
	validateResponse, err := h.s.ValidateDomainChicCRMServices(validateDomainRequest)
	if err != nil {
		switch err.Error() {
		case "username must be a valid email address", "only company email allowed", "username already exists":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error(), "match": false})
		case "domain does not match. To proceed, please check your email":
			c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error(), "match": false})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error(), "match": false})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Domain matches", "match": true, "data": validateResponse})
}
