package ginutils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
)

func ErrHook(c *gin.Context, e error) (int, interface{}) {

	errcode, errpl := 500, e.Error()
	if _, ok := e.(tonic.BindError); ok {
		errcode, errpl = 400, e.Error()
	} else {
		switch {
		case errors.Is(e, errors.BadRequest) || errors.Is(e, errors.NotValid) || errors.Is(e, errors.AlreadyExists) || errors.Is(e, errors.NotSupported) || errors.Is(e, errors.NotAssigned) || errors.Is(e, errors.NotProvisioned):
			errcode, errpl = http.StatusBadRequest, e.Error()
		case errors.Is(e, errors.Forbidden):
			errcode, errpl = http.StatusForbidden, e.Error()
		case errors.Is(e, errors.MethodNotAllowed):
			errcode, errpl = http.StatusMethodNotAllowed, e.Error()
		case errors.Is(e, errors.NotFound) || errors.Is(e, errors.UserNotFound):
			errcode, errpl = http.StatusNotFound, e.Error()
		case errors.Is(e, errors.Unauthorized):
			errcode, errpl = http.StatusUnauthorized, e.Error()
		case errors.Is(e, errors.NotImplemented):
			errcode, errpl = http.StatusNotImplemented, e.Error()
		}
	}

	return errcode, gin.H{"error": errpl}
}
