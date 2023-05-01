package main

import (
	"example.com/go-boot/platform/actuator"
	. "example.com/go-boot/platform/config"
	. "example.com/go-boot/platform/initializer"
	"example.com/go-boot/web/app/graph"
	"example.com/go-boot/web/app/oidc"
	"example.com/go-boot/web/app/openapi"
	"example.com/go-boot/web/app/restapi"
	"example.com/go-boot/web/app/sse"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	Router.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	openapi.NewRouter(Router.Group("/openapi"))
	restapi.Routes(Router.Group("/restapi"))
	graph.Routes(Router.Group("/graphql"))
	oidc.Routes(Router.Group("/login"))
	sse.Routes(Router.Group("/sse"))
	actuator.Routes(Router.Group("/actuator"))

	//Todo set Host - only for local test
	Router.Run(":" + AppConfig.Server.PortNumber)
}
