package utils

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManifestNotFound(t *testing.T) {
	version := GetManifestVersion()

	assert.Equal(t, "Unknown", version, "If no manifest file exist, it should return unknown")

	man = nil
}

func TestManifestFound(t *testing.T) {
	manifestPath = "manifest_test.json"
	data := manifest{
		Name:    "Otto",
		Version: "v0.1.0",
	}

	file, _ := json.MarshalIndent(data, "", " ")
	_ = os.WriteFile(manifestPath, file, 0644)

	version := GetManifestVersion()
	assert.Equal(t, data.Version, version, "Since the file is created the version should match")

	_ = os.Remove(manifestPath)
}
