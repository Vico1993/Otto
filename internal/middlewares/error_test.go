package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIsKnowErrorWithNotNoneError(t *testing.T) {
	assert.Nil(t, isKnowError("foo"), "The error foo doesn't exists")
}

func TestIsKnowErrorWithError(t *testing.T) {
	assert.Equal(t, INTERNAL_ERROR, isKnowError(INTERNAL_ERROR.Message), "The error INTERNAL_ERROR.Message should exists")
}

func TestErrToStatusCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.Use(Error())
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New(BAD_REQUEST.Message))
	})

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode, "Invalid the status code "+strconv.Itoa(http.StatusBadRequest))
}
