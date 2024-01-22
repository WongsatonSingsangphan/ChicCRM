package mail

import (
	"chicCRM/modules/users/mail/handlers"
	"chicCRM/modules/users/mail/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutesMail(router *gin.Engine) {

	s := services.NewServiceAdapter()
	h := handlers.NewHanerhandlerAdapter(s)

	router.POST("/api/mailChicCRM", h.MailChicCRMHandlers)
}
