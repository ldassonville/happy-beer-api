package dispensers

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/beer-puller-api/pkg/api"
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

func (h *Handler) Get(c *gin.Context) {

	ref := c.Param("ref")

	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Access-Control-Allow-Headers, Access-Control-Allow-Origin")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH")

	dispenser, err := h.svc.GetByRef(c.Request.Context(), ref)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "dispenser not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, dispenser)

}

func (h *Handler) Search(c *gin.Context) {

	dispensers, err := h.svc.Search(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if dispensers == nil {
		dispensers = []*api.Dispenser{}
	}

	c.IndentedJSON(http.StatusOK, dispensers)
}

func (h *Handler) Delete(c *gin.Context) {

	ref := c.Param("ref")

	err := h.svc.DeleteByRef(c.Request.Context(), ref)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusOK)
	}
}

func (h *Handler) Create(c *gin.Context) {

	beerDispenser := &api.DispenserEditable{}

	if err := c.ShouldBindJSON(&beerDispenser); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.Create(c.Request.Context(), beerDispenser)
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
