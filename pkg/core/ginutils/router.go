package ginutils

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	logger "github.com/ldassonville/beer-puller-api/pkg/core/logging"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func NewEngine() (*gin.Engine, error) {

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gintrace.Middleware("happy-beer-api"))
	engine.Use(logger.HttpLoggingMiddleware)
	engine.Use(cors.Default())

	return engine, nil
}
