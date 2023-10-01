package handlers

import (
	"encoding/json"
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
