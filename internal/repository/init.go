package repository

import (
	"fmt"
)

var Chat IChatRepository

func Init() {
	Chat = newChatRepository()

	fmt.Println("Repository Initiated")
}
