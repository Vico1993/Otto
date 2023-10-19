package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Vico1993/Otto/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
		Version string `json:"version"`
	}

	utils.ManifestFilePath = "manifestFound.json"
	file, _ := json.MarshalIndent(gin.H{
		"Name":    "Otto",
		"Version": "v0.1.0",
	}, "", " ")
	_ = os.WriteFile(utils.ManifestFilePath, file, 0644)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	Ping(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.Equal(t, "v0.1.0", res.Version)

	_ = os.Remove(utils.ManifestFilePath)
}
