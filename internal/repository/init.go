package repository

import (
	"fmt"
)

var Chat IChatRepository
var Article IArticleRepository

func Init() {
	Chat = newChatRepository()
	Article = newArticleRepository()

	fmt.Println("Repository Initiated")
}
