package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	manifestFilePath = "manifest.json"
	versionNotFound  = "Wasn't able to retrieve versions"
)

// Enpoint to make sure the API is alive
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"version": retrieveVersion(),
	})
}

type manifest struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

func retrieveVersion() string {
	manifestFile, err := os.Open(manifestFilePath)
	if err != nil {
		fmt.Println(err)
		return versionNotFound
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer manifestFile.Close()

	byte, _ := io.ReadAll(manifestFile)

	var manifest manifest

	err = json.Unmarshal(byte, &manifest)
	if err != nil || manifest.Version == "" {
		fmt.Println("Erorr parsign json data: ", err)
		return versionNotFound
	}

	return manifest.Version
}
