package cron

import (
	"testing"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCategoriesAndTagsMatch(t *testing.T) {
	// Override the tags list for the test
	match := isCategoriesAndTagsMatch([]string{"cat1", "cat2", "tag1"}, []string{"tag1", "tag2"})

	assert.Len(t, match, 1, "Tag1 is in the list of category and tag so it should return true")
	assert.Equal(t, []string{"tag1"}, match, "Tag1 is in the list of category and tag so it should return true")
}

func TestCategoriesAndTagsWontMatch(t *testing.T) {
	// Override the tags list for the test
	match := isCategoriesAndTagsMatch([]string{"cat1", "cat2"}, []string{"tag1", "tag2"})

	assert.Len(t, match, 0, "The list of tags and categories don't overlap, array should be empty")
	assert.Equal(t, []string{}, match, "The list of tags and categories don't overlap, array should be empty")
}

type mockArticlRep struct {
	mock.Mock
}

func (m mockArticlRep) Create(title string, published string, link string, source string, tags ...string) *database.Article {
	return nil
}

func (m mockArticlRep) Find(key string, val string) *database.Article {
	return nil
}

func TestFetchFeed(t *testing.T) {
	repository.Article = new(mockArticlRep)
}
