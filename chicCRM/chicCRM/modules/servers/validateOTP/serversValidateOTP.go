package validateOTP

import (
	"chicCRM/pkg/auth/validateOTP/handlers"
	"chicCRM/pkg/auth/validateOTP/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutesValidateOTP(router *gin.Engine) {

	s := services.NewServiceAdapter()
	h := handlers.NewHanerhandlerAdapter(s)

	router.POST("/api/sendOTPEmail", h.RequestEmailForValidateOTPChicCRMHandlers)
	router.POST("/api/validateOTPEmail", h.ValidateOTPFromRequestEmailChicCRMHandlers)
	router.POST("/api/qrTOTP", h.QrTOTPChicCRMHandlers)
	router.POST("/api/validateQrTOTP", h.ValidateQrTOTPChicCRMHandlers)
	router.DELETE("/api/deleteKeyQrTOTP", h.DeleteKeyQrTOTPChicCRMHandlers)
}
