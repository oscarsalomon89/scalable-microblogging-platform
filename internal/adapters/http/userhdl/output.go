package userhdl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pkgerrors "github.com/oscarsalomon89/go-hexagonal/pkg/errors"
)

// Centraliza la respuesta de errores para los handlers
func handleError(c *gin.Context, err error) {
	switch {
	case pkgerrors.IsValidationError(err):
		c.JSON(http.StatusBadRequest, errorResponse(err))
	case pkgerrors.IsNotFound(err):
		c.JSON(http.StatusNotFound, errorResponse(err))
	case pkgerrors.IsConflict(err):
		c.JSON(http.StatusConflict, errorResponse(err))
	case pkgerrors.IsAuthenticationError(err):
		c.JSON(http.StatusUnauthorized, errorResponse(err))
	case pkgerrors.IsAuthorizationError(err):
		c.JSON(http.StatusForbidden, errorResponse(err))
	default:
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}
}

func errorResponse(err error) gin.H {
	if e, ok := err.(*pkgerrors.Error); ok {
		return gin.H{"error": e.ToJSON()}
	}
	return gin.H{"error": err.Error()}
}
