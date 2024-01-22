package handlers

import (
	"chicCRM/pkg/auth/validateOTP/models"
	"chicCRM/pkg/auth/validateOTP/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPort interface {
	RequestEmailForValidateOTPChicCRMHandlers(c *gin.Context)
	ValidateOTPFromRequestEmailChicCRMHandlers(c *gin.Context)
	QrTOTPChicCRMHandlers(c *gin.Context)
	ValidateQrTOTPChicCRMHandlers(c *gin.Context)
	DeleteKeyQrTOTPChicCRMHandlers(c *gin.Context)
}

type handlerAdapter struct {
	s services.ServicePort
}

func NewHanerhandlerAdapter(s services.ServicePort) HandlerPort {
	return &handlerAdapter{s: s}
}

func (h *handlerAdapter) RequestEmailForValidateOTPChicCRMHandlers(c *gin.Context) {

	var requestBody map[string]string
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	email := requestBody["email"]
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": "Please provide your email"})
		return
	}
	referenceID, err := h.s.RequestEmailForValidateOTPChicCRMServices(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "The OTP has been sent to your email.", "referenceID": referenceID})
}

func (h *handlerAdapter) ValidateOTPFromRequestEmailChicCRMHandlers(c *gin.Context) {
	var validateBody models.ValidateBody
	// fmt.Println(validateBody.Email, validateBody.OTP)
	if err := c.ShouldBindJSON(&validateBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	// err := h.s.ValidateOTPFromRequestEmailChicCRMServices(validateBody.Email, validateBody.OTP, validateBody.ReferenceID)
	err := h.s.ValidateOTPFromRequestEmailChicCRMServices(validateBody.OTP, validateBody.ReferenceID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Valid OTP"})
}

func (h *handlerAdapter) QrTOTPChicCRMHandlers(c *gin.Context) {
	var qrTOTPRequest models.QrTOTPRequest
	if err := c.ShouldBindJSON(&qrTOTPRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	qrCodeURL, statusqr, err := h.s.QrTOTPChicCRMServices(qrTOTPRequest.AccountName, qrTOTPRequest.Value)
	if err != nil {
		switch err.Error() {
		case "AccountName already exists":
			c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error(), "statusqr": statusqr})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": err.Error()})
		}
		return
	}
	if statusqr {
		c.JSON(http.StatusOK, gin.H{"qrCodeURL": qrCodeURL, "statusqr": statusqr, "status": "OK"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please enter 1 to process", "status": "Error"})
	}
}

func (h *handlerAdapter) ValidateQrTOTPChicCRMHandlers(c *gin.Context) {
	var validateQrTOTP models.ValidateQrTOTP
	if err := c.ShouldBindJSON(&validateQrTOTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	isValid, err := h.s.ValidateQrTOTPChicCRMServices(validateQrTOTP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": "Error"})
		return
	}

	if isValid {
		c.JSON(http.StatusOK, gin.H{"message": "OTP is valid", "status": "OK"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "OTP is invalid", "status": "Error"})
	}
}

func (h *handlerAdapter) DeleteKeyQrTOTPChicCRMHandlers(c *gin.Context) {
	var deleteKeyQrTOTP models.DeleteKeyQrTOTP
	if err := c.ShouldBindJSON(&deleteKeyQrTOTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Error", "message": err.Error()})
		return
	}
	err := h.s.DeleteKeyQrTOTPChicCRMServices(deleteKeyQrTOTP.AccountName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Error", "message": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "OTP key for " + deleteKeyQrTOTP.AccountName + " has been deleted"})
	}
}
