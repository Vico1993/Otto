package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/subosito/gotenv"
)

type initData struct {
	ChatId string   `json:"chatid"`
	UserId string   `json:"userid"`
	Tags   []string `json:"tags"`
	Feeds  []string `json:"feeds"`
}

func main() {
	jsonFile, err := os.Open("./scripts/init-data.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Opened init-data.json")
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)

	var data initData

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatal(err)
	}

	userId, err := strconv.ParseInt(data.UserId, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// load .env file if any otherwise use env set
	fmt.Println("Loading the env variables")
	_ = gotenv.Load()

	fmt.Println("Initialization of the database connection")
	database.Init()
	repository.Init()

	fmt.Println("Pushing in Database")
	repository.Chat.Create(
		data.ChatId,
		userId,
		data.Tags,
		data.Feeds,
	)

	fmt.Println("Insert done!!")
}
