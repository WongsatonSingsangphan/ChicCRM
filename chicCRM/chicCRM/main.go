package main

import (
	"chicCRM/modules/servers/finalCode"
	"chicCRM/modules/servers/login"
	"chicCRM/modules/servers/mail"
	"chicCRM/modules/servers/migrateData"
	"chicCRM/modules/servers/password"
	"chicCRM/modules/servers/queryData"
	"chicCRM/modules/servers/register"
	"chicCRM/modules/servers/uploadBinary"
	"chicCRM/modules/servers/validateOTP"
	"chicCRM/pkg/databases/postgresql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := postgresql.Postgresql()
	defer db.Close()

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "X-Auth-Token", "Authorization"}
	router.Use(cors.New(config))

	queryData.SetupRoutesQueryData(router, db)
	register.SetupRoutesRegister(router, db)
	finalCode.FinalCode(router, db)
	login.SetupRoutesLogin(router, db)
	password.SetupRoutesInitPassword(router, db)
	mail.SetupRoutesMail(router)
	validateOTP.SetupRoutesValidateOTP(router)
	migrateData.SetupRoutesMigrateDataByOrganize(router, db)
	uploadBinary.SetupRoutesUploadBinary(router, db)

	err := router.Run(":8888")
	if err != nil {
		panic(err.Error())
	}
}
