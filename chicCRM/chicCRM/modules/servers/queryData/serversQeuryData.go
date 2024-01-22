package queryData

import (
	"chicCRM/modules/queryData/handlers"
	"chicCRM/modules/queryData/repositories"
	"chicCRM/modules/queryData/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoutesQueryData(router *gin.Engine, db *sql.DB) {

	r := repositories.NewRepositoryAdapter(db)
	s := services.NewServiceAdapter(r)
	h := handlers.NewHanerhandlerAdapter(s)

	router.GET("/api/getCountries", h.GetCountriesHandlers)
	router.GET("/api/getProvinceAmphoeTambonZipcode", h.GetProvinceAmphoeTambonZipcodeHandlers)
	router.GET("/api/getUniversalInfo", h.GetUniversalInfoHandlers)
}
