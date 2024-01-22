package finalCode

import (
	"chicCRM/modules/finalCode/handlers"
	"chicCRM/modules/finalCode/repositories"
	"chicCRM/modules/finalCode/services"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func FinalCode(router *gin.Engine, db *sql.DB) {

	r := repositories.NewRepositoryAdapter(db)
	s := services.NewServiceAdapter(r)
	h := handlers.NewHanerhandlerAdapter(s, db)

	router.POST("/api/requestDoc", h.AddActivity)
	router.GET("/api/requestFile/:scdact_id", h.GETUser_File)
	router.POST("/api/fileEncrypt", h.File_encrypt)
	router.POST("/api/sendMailFinalcode", h.Send_mail)
	router.GET("/api/checkOrganizeFeature/:id", h.CheckFeatureOrganizeHandlers)
	router.GET("/api/checkMemberFeature/:id", h.CheckFeatureMemberHandlers)
	router.GET("/api/getPolicyData/:id", h.GetPolicyDataHandlers)
	router.GET("/api/getLogSecuredocActivityByMember/:uuid", h.GetSecuredocActivityByOrganizeMemberUUIDHandlers)
	router.GET("/api/getlogSecuredocActivityByTeamlead/:uuid", h.GetSecuredocActivityByTeamleadIDHandlers)
	router.GET("/api/getPolicyAuthorizationByTeamlead/:uuid", h.GetPolicyAuthorizationByTeamleadHandlers)
	router.GET("/api/getStatusByRequestID/:requestID", h.GetRequestStatusByReqIDHandlers)
}
