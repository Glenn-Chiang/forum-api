package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Sends a error response in JSON with status code and message
func HTTPErrorResponse(ctx *gin.Context, err error) {
	var e *Error
	// If the error is not one of the defined custom errors, return a generic internal server error
	if !errors.As(err, &e) {
		ctx.JSON(http.StatusInternalServerError, errorResponseBody(err))
	}
	// Choose appropriate status code based on custom error code
	switch e.Code {
	case ErrInvalid:
		ctx.JSON(http.StatusBadRequest, errorResponseBody(e))
	case ErrUnauthorized:
		ctx.JSON(http.StatusUnauthorized, errorResponseBody(e))
	case ErrNotFound:
		ctx.JSON(http.StatusNotFound, errorResponseBody(e))
	case ErrConflict:
		ctx.JSON(http.StatusConflict, errorResponseBody(e))
	case ErrInternal:
		ctx.JSON(http.StatusInternalServerError, errorResponseBody(e))
	}
}

func errorResponseBody(err error) any {
	return gin.H{"error": err.Error()}
}
