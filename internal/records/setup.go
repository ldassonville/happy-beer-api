package records

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/beer-puller-api/internal/records/storage"
	"github.com/ldassonville/beer-puller-api/pkg/core/ioc"
)

const (
	IocRecordSvc = "records-service"
)

func Setup(injector *ioc.Injector) {

	engine := injector.Get(ioc.IocEngine).(*gin.Engine)

	recordDao := storage.NewMemoryDao()
	recordSvc := NewService(recordDao)
	injector.Register(IocRecordSvc, recordSvc)

	handler := NewHandler(recordSvc)

	engine.Handle("POST", "/records", handler.Create)
	engine.Handle("GET", "/records", handler.Search)
	engine.Handle("GET", "/records/_all/events", handler.Events)
}
