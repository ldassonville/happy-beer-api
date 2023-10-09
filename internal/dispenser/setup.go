package dispenser

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/beer-puller-api/internal/dispenser/storage"
	"github.com/ldassonville/beer-puller-api/pkg/core/ioc"
)

func Setup(injector *ioc.Injector) {

	engine := injector.Get(ioc.IocEngine).(*gin.Engine)

	dispenserDao := storage.NewMemoryDao()
	productSvc := &Service{
		dao: dispenserDao,
	}

	handler := NewHandler(productSvc)
	engine.Handle("POST", "/dispensers", handler.Create)
	engine.Handle("GET", "/dispensers/:ref", handler.Get)
	engine.Handle("GET", "/dispensers", handler.Search)
	engine.Handle("DELETE", "/dispensers", handler.Delete)
}
