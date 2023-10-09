package dispenser

import (
	"github.com/gin-gonic/gin"
	"github.com/ldassonville/beer-puller-api/pkg/model"
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
	c.IndentedJSON(http.StatusOK, dispensers)
}

func (h *Handler) Delete(c *gin.Context) {

	ref := c.Param("ref")

	err := h.svc.DeleteByRef(c.Request.Context(), ref)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) Create(c *gin.Context) {

	beerDispenser := &model.Dispenser{}

	if err := c.BindJSON(&beerDispenser); err != nil {
		return
	}

	h.svc.Create(c.Request.Context(), beerDispenser)

	c.IndentedJSON(http.StatusCreated, beerDispenser)
}
