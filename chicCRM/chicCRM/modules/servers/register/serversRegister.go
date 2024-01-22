package register

import (
	"chicCRM/modules/users/register/handlers"
	"chicCRM/modules/users/register/repositories"
	"chicCRM/modules/users/register/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutesRegister(router *gin.Engine, db *sql.DB) {

	r := repositories.NewRepositoryAdapter(db)
	s := services.NewServiceAdapter(r)
	h := handlers.NewHanerhandlerAdapter(s)

	router.POST("/api/registerChicCRM", h.RegisterChicCRMHandlers)
	router.POST("/api/validateDomainChicCRM", h.ValidateDomainChicCRMHandlers)
}
