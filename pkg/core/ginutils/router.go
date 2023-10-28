package ginutils

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewEngine() (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	engine.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:9000", "http://localhost:4200"}
	//config.AllowAllOrigins = true

	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "X-Api-Key", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin"}
	config.AllowMethods = []string{"POST", "PUT", "PATCH", "GET", "DELETE", "OPTIONS"}
	engine.Use(cors.New(config))

	//Access-Control-Allow-Origin: http://example.com

	//engine.Use(gintrace.Middleware("happy-beer-api"))
	//engine.Use(logger.HttpLoggingMiddleware)

	engine.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return engine, nil
}
