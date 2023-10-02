package repository

import "fmt"

var (
	Chat    IChatRepository
	Feed    IFeedRepository
	Article IArticleRepository
)

func Init() {
	Chat = &SChatRepository{}
	Feed = &SFeedRepository{}
	Article = &SArticleRepository{}

	fmt.Println("Repository Initiated")
}
