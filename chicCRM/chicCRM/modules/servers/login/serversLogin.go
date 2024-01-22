package login

import (
	"chicCRM/modules/users/login/handlers"
	"chicCRM/modules/users/login/repositories"
	"chicCRM/modules/users/login/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutesLogin(router *gin.Engine, db *sql.DB) {

	r := repositories.NewRepositoryAdapter(db)
	s := services.NewServiceAdapter(r)
	h := handlers.NewHanerhandlerAdapter(s)

	router.POST("/api/LoginChicCRM", h.LoginChicCRMHandlers)
	router.POST("/api/LoginTeamleadSecuredoc", h.LoginTeamleadSecuredocHandlers)
}
