package records

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/happy-beer-api/pkg/api"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Search(c *gin.Context) {

	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")

	events, err := h.svc.Search(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if events == nil {
		events = []*api.Record{}
	}

	c.IndentedJSON(http.StatusOK, events)
}

func (h *Handler) Create(c *gin.Context) {

	event := &api.Record{}

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Create(c.Request.Context(), event)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func (h *Handler) Events(c *gin.Context) {

	sub := h.svc.Subscribe()

	c.Stream(func(w io.Writer) bool {
		if event, ok := <-sub.Chan(); ok {
			c.SSEvent(event.Typ, event.Msg)
			return true
		}
		sub.Leave()
		return false
	})

	logrus.Infof("SSE subscription %s finish", sub.Id())

}
