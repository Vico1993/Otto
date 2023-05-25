package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoriesAndTagsMatch(t *testing.T) {
	// Override the tags list for the test
	tags = []string{"tag1", "tag2"}
	match := isCategoriesAndTagsMatch([]string{"cat1", "cat2", "tag1"})

	assert.Len(t, match, 1, "Tag1 is in the list of category and tag so it should return true")
	assert.Equal(t, []string{"tag1"}, match, "Tag1 is in the list of category and tag so it should return true")
}

func TestCategoriesAndTagsWontMatch(t *testing.T) {
	// Override the tags list for the test
	tags = []string{"tag1", "tag2"}
	match := isCategoriesAndTagsMatch([]string{"cat1", "cat2"})

	assert.Len(t, match, 0, "The list of tags and categories don't overlap, array should be empty")
	assert.Equal(t, []string{}, match, "The list of tags and categories don't overlap, array should be empty")
}
