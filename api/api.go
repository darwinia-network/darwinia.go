// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/darwinia-network/shadow/api/docs"
	"github.com/darwinia-network/shadow/internal/util"
)

func Swagger() {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := gin.Default()
	c, err := NewController()
	util.Assert(err)

	v1 := r.Group("/api/v1")
	{
		header := v1.Group("/header")
		{
			header.GET("/hash/:hash", c.GetEthHeaderByHash)
			header.GET("/number/:number", c.GetEthHeaderByNumber)
		}
		proof := v1.Group("/proof")
		{
			proof.GET("/hash/:hash", c.GetEthProofByHash)
			proof.GET("/number/:number", c.GetEthProofByNumber)
		}
		header.POST("/proposal", c.GetEthHeaderByHash)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	gin.SetMode(gin.ReleaseMode)
	util.Assert(r.Run(":3600"))
}
