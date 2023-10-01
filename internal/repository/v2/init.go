package v2

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
