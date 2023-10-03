package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveVersionFileNotFound(t *testing.T) {
	manifestFilePath = "manifestNotFound.json"

	res := retrieveVersion()
	assert.Equal(t, versionNotFound, res, "Looking for a file that doesn't exist, should return default")
}

func TestRetrieveVersionFileFound(t *testing.T) {
	manifestFilePath = "manifestFound.json"
	file, _ := json.MarshalIndent(gin.H{
		"Name":    "Otto",
		"Version": "v0.1.0",
	}, "", " ")
	_ = os.WriteFile(manifestFilePath, file, 0644)

	res := retrieveVersion()
	assert.Equal(t, "v0.1.0", res, "Looking for a file that doesn't exist, should return default")

	_ = os.Remove(manifestFilePath)
}

func TestRetrieveVersionFileFoundIncorrectJson(t *testing.T) {
	manifestFilePath = "manifestFound.json"
	file, _ := json.MarshalIndent(gin.H{
		"Name":     "Otto",
		"Versions": "v0.1.0",
	}, "", " ")
	_ = os.WriteFile(manifestFilePath, file, 0644)

	res := retrieveVersion()
	assert.Equal(t, versionNotFound, res, "File is created, but json incorrect. Should return default version not found message")

	_ = os.Remove(manifestFilePath)
}

func TestPing(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
		Version string `json:"version"`
	}

	manifestFilePath = "manifestFound.json"
	file, _ := json.MarshalIndent(gin.H{
		"Name":    "Otto",
		"Version": "v0.1.0",
	}, "", " ")
	_ = os.WriteFile(manifestFilePath, file, 0644)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	Ping(ctx)

	var res Response
	_ = json.Unmarshal(recorder.Body.Bytes(), &res)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode, "StatusCode should be OK")
	assert.Equal(t, "v0.1.0", res.Version)

	_ = os.Remove(manifestFilePath)
}
