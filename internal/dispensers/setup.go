package dispensers

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/happy-beer-api/internal/dispensers/storage"
	"github.com/ldassonville/happy-beer-api/internal/records"
	"github.com/ldassonville/happy-beer-api/pkg/core/ioc"
)

func Setup(injector *ioc.Injector) {

	engine := injector.Get(ioc.IocEngine).(*gin.Engine)
	recordSvc := injector.Get(records.IocRecordSvc).(*records.Service)

	dispenserDao := storage.NewMemoryDao()
	dispenserSvc := NewService(dispenserDao, recordSvc)

	handler := NewHandler(dispenserSvc)
	engine.Handle("POST", "/dispensers", handler.Create)
	engine.Handle("GET", "/dispensers/:ref", handler.Get)
	engine.Handle("GET", "/dispensers", handler.Search)
	engine.Handle("DELETE", "/dispensers/:ref", handler.Delete)
	engine.Handle("PUT", "/dispensers/:ref", handler.Update)

	engine.Handle("GET", "/dispensers/_all/events", handler.Events)
}
