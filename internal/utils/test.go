package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MockPostRequest(c *gin.Context, content interface{}, isPut bool) {
	if c.Request == nil {
		c.Request = &http.Request{
			Header: make(http.Header),
		}
	}

	if isPut {
		c.Request.Method = "PUT"
	} else {
		c.Request.Method = "POST"
	}

	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
