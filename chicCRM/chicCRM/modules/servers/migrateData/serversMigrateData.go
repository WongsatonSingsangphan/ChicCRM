package migrateData

import (
	"chicCRM/modules/migrateData/handlers"
	"chicCRM/modules/migrateData/repositories"
	"chicCRM/modules/migrateData/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutesMigrateDataByOrganize(router *gin.Engine, db *sql.DB) {

	r := repositories.NewRepositoryAdapter(db)
	s := services.NewServiceAdapter(r)
	h := handlers.NewHanerhandlerAdapter(s)

	router.POST("/api/migrateDataByOrganize", h.MigrateMemberDataByOrganizeHandlers)
	router.POST("/api/migrateTeamleadDataByOrganize", h.MigrateTeamleadTracByOrganizeHandlers)
}
