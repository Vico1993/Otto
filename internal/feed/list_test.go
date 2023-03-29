package feed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetList(t *testing.T) {
	// Override the tags list for the test
	tags = []string{"tag1"}

	// Override the tags list for the test
	listOfFeeds = []string{"https://google.com/feed/"}

	urls := GetList()

	// assert.Len(
	// 	t,
	// 	urls,
	// 	2,
	// 	"getList should return 2 elements as tags should be 1 and there is only 1 list",
	// )
	// assert.Equal(
	// 	t,
	// 	[]string{"https://medium.com/feed/tag/tag1", "https://google.com/feed/"},
	// 	urls,
	// 	"Response doesn't match the expected value",
	// )

	assert.Len(
		t,
		urls,
		1,
		"getList should return 1 elements as medium is removed for now",
	)
	assert.Equal(
		t,
		[]string{"https://google.com/feed/"},
		urls,
		"Response doesn't match the expected value",
	)
}
