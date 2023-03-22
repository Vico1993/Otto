package feed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoriesAndTagsMatch(t *testing.T) {
	// Override the tags list for the test
	tags = []string{"tag1", "tag2"}

	assert.True(t, isCategoriesAndTagsMatch([]string{"cat1", "cat2", "tag1"}), "Tag1 is in the list of category and tag so it should return true")
}

func TestCategoriesAndTagsWontMatch(t *testing.T) {
	// Override the tags list for the test
	tags = []string{"tag1", "tag2"}

	assert.False(t, isCategoriesAndTagsMatch([]string{"cat1", "cat2"}), "The list of tags and categories don't overlap")
}
