package dispensers

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/happy-beer-api/pkg/api"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
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

	if dispenser.Status.Status == api.InternalStatusArchived {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "dispenser is archived"})
		return
	}

	c.IndentedJSON(http.StatusOK, dispenser)

}

func (h *Handler) Search(c *gin.Context) {

	query := &api.DispenserQuery{}

	statuz, _ := c.GetQueryArray("internalStatus")
	for _, status := range statuz {
		if strings.EqualFold(status, "all") {
			query.Statuses = append(query.Statuses, api.InternalStatusActive, api.InternalStatusArchived, api.InternalStatusPending)
			break
		} else {
			query.Statuses = append(query.Statuses, api.InternalStatus(status))
		}
	}

	if len(query.Statuses) == 0 {
		query.Statuses = []api.InternalStatus{api.InternalStatusActive}
	}

	dispensers, err := h.svc.Search(c.Request.Context(), query)
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
		c.Status(http.StatusNoContent)
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

func (h *Handler) Update(c *gin.Context) {

	editableDispenser := api.DispenserEditable{}

	if err := c.ShouldBindJSON(&editableDispenser); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	editableDispenser.Ref = c.Param("ref")

	beerDispenser := &api.Dispenser{
		DispenserEditable: editableDispenser,
	}

	res, err := h.svc.Update(c.Request.Context(), beerDispenser)
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
