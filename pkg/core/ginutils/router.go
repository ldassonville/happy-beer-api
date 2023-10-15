package ginutils

import (
	"github.com/gin-gonic/gin"
)

func NewEngine() (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)

	engine := gin.Default()

	engine.Use(func(c *gin.Context) {

		c.Writer.Header().Add("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Access-Control-Allow-Headers, Access-Control-Allow-Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {

			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	//engine.Use(gin.Recovery())
	//engine.Use(gintrace.Middleware("happy-beer-api"))
	//engine.Use(logger.HttpLoggingMiddleware)

	return engine, nil
}
