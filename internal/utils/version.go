package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type manifest struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

var man *manifest = nil
var manifestPath string = "manifest.json"

// Return Manifest version
func GetManifestVersion() string {
	if man == nil {
		jsonFile, err := os.Open(manifestPath)
		if err != nil {
			fmt.Println("Couldn't find the manifest file...")
			man = &manifest{
				Name:    "Otto",
				Version: "Unknown",
			}
		} else {
			defer jsonFile.Close()

			byteValue, _ := io.ReadAll(jsonFile)
			_ = json.Unmarshal([]byte(byteValue), &man)
		}
	}

	return man.Version
}
