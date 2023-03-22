package feed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMediumUrl(t *testing.T) {
	// Override the tags list for the test
	tags = []string{"tag1", "tag2"}

	// Call the buildMediumFeedBasedOnTag() method
	urls := buildMediumFeedBasedOnTag()

	assert.Len(
		t,
		urls,
		2,
		"BuildMediumFeedBasedOnTag should return 2 elements as getInterestTags mock only return 2 elements",
	)
	assert.Equal(
		t,
		[]string{"https://medium.com/feed/tag/tag1", "https://medium.com/feed/tag/tag2"},
		urls,
		"Response doesn't match the expected value",
	)
}
