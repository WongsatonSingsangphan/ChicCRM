package password

import (
	"chicCRM/modules/middlewares"
	"chicCRM/modules/users/Password/handlers"
	"chicCRM/modules/users/Password/repositories"
	"chicCRM/modules/users/Password/services"

	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutesInitPassword(router *gin.Engine, db *sql.DB) {

	r := repositories.NewRepositoryAdapter(db)
	s := services.NewServiceAdapter(r)
	h := handlers.NewHanerhandlerAdapter(s)

	router.PATCH("/api/InitPasswordChicCRM", middlewares.AuthMiddleware(db), h.InitPasswordChicCRMHandlers)
	router.PATCH("/api/ChangePasswordChicCRM", h.ChangePasswordChicCRMHandlers)
	router.POST("/api/RequestResetPasswordChicCRM", h.RequestResetPasswordHandlers)
	router.PATCH("/api/ResetPasswordChicCRM", middlewares.AuthMiddlewareResetPassword(db), h.ResetPasswordChicCRMHandlers)
}
