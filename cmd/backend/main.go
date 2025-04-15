package main

import (
	"backend/internal/application"
	"backend/internal/infrastructure"
)

//	@title			Swagger API
//	@version		1.0
//	@description	T-Loyal API
//	@title			Crazy Peppers T-Loyal API

//	@host		prod-team-29-4254c2ee.REDACTED
//	@BasePath	/api/v1

//	@schemes	http https
//	@produce	application/json
//	@consumes	application/json

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	config := infrastructure.LoadConfig()
	app := application.NewApplication(config)
	if err := app.RunServer(); err != nil {
		return
	}
}
