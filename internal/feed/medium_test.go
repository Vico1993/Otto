package feed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMediumUrl(t *testing.T) {
	// Override the tags list for the test
	tags = []string{"tag1", "tag2"}

	urls := buildMediumFeedBasedOnTag()

	assert.Len(
		t,
		urls,
		2,
		"BuildMediumFeedBasedOnTag should return 2 elements as tags mock only return 2 elements",
	)
	assert.Equal(
		t,
		[]string{"https://medium.com/feed/tag/tag1", "https://medium.com/feed/tag/tag2"},
		urls,
		"Response doesn't match the expected value",
	)
}

func TestBuildMediumUrlWithEmptyTags(t *testing.T) {
	// Override the tags list for the test
	tags = []string{}

	urls := buildMediumFeedBasedOnTag()

	assert.Len(
		t,
		urls,
		0,
		"BuildMediumFeedBasedOnTag should be empty if no tag are set",
	)
	assert.Equal(
		t,
		[]string{},
		urls,
		"Response doesn't match the expected value",
	)
}
