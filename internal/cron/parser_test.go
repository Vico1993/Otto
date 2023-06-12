package cron

import (
	"errors"
	"testing"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestCategoriesAndTagsMatch(t *testing.T) {
	p := &parser{
		url:  "https://test.com",
		tags: []string{"tag1", "tag2"},
	}

	match := p.isCategoriesAndTagsMatch([]string{"cat1", "cat2", "tag1"})

	assert.Len(t, match, 1, "Tag1 is in the list of category and tag so it should return true")
	assert.Equal(t, []string{"tag1"}, match, "Tag1 is in the list of category and tag so it should return true")
}

func TestCategoriesAndTagsWontMatch(t *testing.T) {
	p := &parser{
		url:  "https://test.com",
		tags: []string{"tag1", "tag2"},
	}

	match := p.isCategoriesAndTagsMatch([]string{"cat1", "cat2"})

	assert.Len(t, match, 0, "The list of tags and categories don't overlap, array should be empty")
	assert.Equal(t, []string{}, match, "The list of tags and categories don't overlap, array should be empty")
}

func TestExecuteSuccessfulFeed(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{"tag1", "tag2"},
	}

	item := &gofeed.Item{
		Title:     "Super Title for an Article",
		Published: "2023-06-04",
		Link:      "https://test.com/article-1",
		Categories: []string{
			"tag1", "tag5", "tag8",
		},
		Authors: []*gofeed.Person{
			{
				Name: "Victor",
			},
		},
	}

	f := &gofeed.Feed{
		Title: "Super Test",
		Items: []*gofeed.Item{
			item,
		},
	}

	oldParseUrl := parseUrl
	defer func() { parseUrl = oldParseUrl }()

	parseUrl = func(url string) (*gofeed.Feed, error) {
		return f, nil
	}

	// Mock Article Repository
	articleRepositoryMock := new(repository.MocksArticleRep)

	// Article repository return nil on Find
	articleRepositoryMock.On("Find", "title", item.Title).Return(nil)

	articleExpected := database.NewArticle(
		item.Title,
		item.Published,
		item.Link,
		f.Title,
		item.Authors[0].Name,
		[]string{"tag1"},
		item.Categories...,
	)

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		item.Title,
		item.Published,
		item.Link,
		f.Title,
		item.Authors[0].Name,
		[]string{"tag1"},
		item.Categories,
	).Return(articleExpected)

	result, err := p.execute(articleRepositoryMock)

	articleRepositoryMock.AssertCalled(t, "Find", "title", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", item.Title, item.Published, item.Link, f.Title, item.Authors[0].Name, []string{"tag1"}, item.Categories)

	assert.Nil(t, err, "The error object should be nil")
	assert.Len(t, result.articles, 1, "After execute should receive 1 article")

	assert.EqualValues(t, articleExpected.Title, result.articles[0].Title, "The expected article should be the one set on the mock return")
}

func TestExecuteArticleAlreadyInDB(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{"tag1", "tag2"},
	}

	item := &gofeed.Item{
		Title:     "Super Title for an Article",
		Published: "2023-06-04",
		Link:      "https://test.com/article-1",
		Categories: []string{
			"tag1", "tag5", "tag8",
		},
		Authors: []*gofeed.Person{
			{
				Name: "Victor",
			},
		},
	}

	f := &gofeed.Feed{
		Title: "Super Test",
		Items: []*gofeed.Item{
			item,
		},
	}

	oldParseUrl := parseUrl
	defer func() { parseUrl = oldParseUrl }()

	parseUrl = func(url string) (*gofeed.Feed, error) {
		return f, nil
	}

	// Mock Article Repository
	articleRepositoryMock := new(repository.MocksArticleRep)

	articleExpected := database.NewArticle(
		item.Title,
		item.Published,
		item.Link,
		f.Title,
		item.Authors[0].Name,
		[]string{"tag1"},
		item.Categories...,
	)

	// Article repository return article on Find
	articleRepositoryMock.On("Find", "title", item.Title).Return(articleExpected)

	result, err := p.execute(articleRepositoryMock)

	articleRepositoryMock.AssertCalled(t, "Find", "title", item.Title)
	articleRepositoryMock.AssertNotCalled(t, "Create")

	assert.Nil(t, err, "The error object should be nil")
	assert.Len(t, result.articles, 0, "After execute should receive 0 article has it's already found")
}

func TestExecuteArticleNoCategoryMatch(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{"tag1", "tag2"},
	}

	item := &gofeed.Item{
		Title:     "Super Title for an Article",
		Published: "2023-06-04",
		Link:      "https://test.com/article-1",
		Categories: []string{
			"tag5", "tag8",
		},
		Authors: []*gofeed.Person{
			{
				Name: "Victor",
			},
		},
	}

	f := &gofeed.Feed{
		Title: "Super Test",
		Items: []*gofeed.Item{
			item,
		},
	}

	oldParseUrl := parseUrl
	defer func() { parseUrl = oldParseUrl }()

	parseUrl = func(url string) (*gofeed.Feed, error) {
		return f, nil
	}

	// Mock Article Repository
	articleRepositoryMock := new(repository.MocksArticleRep)

	result, err := p.execute(articleRepositoryMock)

	articleRepositoryMock.AssertNotCalled(t, "Find")
	articleRepositoryMock.AssertNotCalled(t, "Create")

	assert.Nil(t, err, "The error object should be nil")
	assert.Len(t, result.articles, 0, "Since no category match, should return 0")
}

func TestExecuteParsingFailed(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{"tag1", "tag2"},
	}

	oldParseUrl := parseUrl
	defer func() { parseUrl = oldParseUrl }()

	parseUrl = func(url string) (*gofeed.Feed, error) {
		return nil, errors.New("Failling...")
	}

	// Mock Article Repository
	articleRepositoryMock := new(repository.MocksArticleRep)

	result, err := p.execute(articleRepositoryMock)

	articleRepositoryMock.AssertNotCalled(t, "Find")
	articleRepositoryMock.AssertNotCalled(t, "Create")

	assert.Nil(t, result, "Result should be nil if gofeed fail")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Couldn't parsed test.com: Failling..."), err, "Error message doesn't match")
	}
}
