package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	ManifestFilePath = "manifest.json"
	VersionNotFound  = "Wasn't able to retrieve versions"
)

type manifest struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

func RetrieveVersion() string {
	manifestFile, err := os.Open(ManifestFilePath)
	if err != nil {
		fmt.Println(err)
		return VersionNotFound
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer manifestFile.Close()

	byte, _ := io.ReadAll(manifestFile)

	var manifest manifest

	err = json.Unmarshal(byte, &manifest)
	if err != nil || manifest.Version == "" {
		fmt.Println("Erorr parsign json data: ", err)
		return VersionNotFound
	}

	return manifest.Version
}
