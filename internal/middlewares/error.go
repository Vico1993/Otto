package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	BAD_REQUEST = &sError{
		Code:    http.StatusBadRequest, // 400
		Message: "Bad Request",
	}
	UNAUTHORIZED = &sError{
		Code:    http.StatusUnauthorized, // 401
		Message: "Unauthorized",
	}
	FORBIDDEN = &sError{
		Code:    http.StatusForbidden, // 403
		Message: "Forbidden",
	}
	NOT_FOUND = &sError{
		Code:    http.StatusNotFound, // 404
		Message: "Not Found",
	}
	METHOD_NOT_ALLOWED = &sError{
		Code:    http.StatusMethodNotAllowed, // 405
		Message: "Method Not Allowed",
	}
	INTERNAL_ERROR = &sError{
		Code:    http.StatusInternalServerError, // 500
		Message: "Internal Server Error",
	}
	list = []*sError{
		INTERNAL_ERROR, METHOD_NOT_ALLOWED, NOT_FOUND, FORBIDDEN, UNAUTHORIZED, BAD_REQUEST,
	}
)

type sError struct {
	Code    int
	Message string
}

// Error Mildleware to format my error
func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastError := c.Errors.Last()
		if lastError == nil {
			return
		}

		err := isKnowError(lastError.Error())
		if err == nil {
			err = INTERNAL_ERROR
		}

		c.JSON(err.Code, gin.H{
			"error": err.Message,
		})
	}
}

// Is error is known
func isKnowError(errStr string) *sError {
	for _, err := range list {
		if errStr == err.Message {
			return err
		}
	}

	return nil
}
