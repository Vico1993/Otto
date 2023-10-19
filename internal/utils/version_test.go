package utils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveVersionFileNotFound(t *testing.T) {
	ManifestFilePath = "manifestNotFound.json"

	res := RetrieveVersion()
	assert.Equal(t, VersionNotFound, res, "Looking for a file that doesn't exist, should return default")
}

func TestRetrieveVersionFileFound(t *testing.T) {
	ManifestFilePath = "manifestFound.json"
	file, _ := json.MarshalIndent(gin.H{
		"Name":    "Otto",
		"Version": "v0.1.0",
	}, "", " ")
	_ = os.WriteFile(ManifestFilePath, file, 0644)

	res := RetrieveVersion()
	assert.Equal(t, "v0.1.0", res, "Looking for a file that doesn't exist, should return default")

	_ = os.Remove(ManifestFilePath)
}

func TestRetrieveVersionFileFoundIncorrectJson(t *testing.T) {
	ManifestFilePath = "manifestFound.json"
	file, _ := json.MarshalIndent(gin.H{
		"Name":     "Otto",
		"Versions": "v0.1.0",
	}, "", " ")
	_ = os.WriteFile(ManifestFilePath, file, 0644)

	res := RetrieveVersion()
	assert.Equal(t, VersionNotFound, res, "File is created, but json incorrect. Should return default version not found message")

	_ = os.Remove(ManifestFilePath)
}
