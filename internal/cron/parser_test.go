package cron

import (
	"testing"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

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
	articleRepositoryMock := new(repository.MocksArticleRepository)

	// Article repository return nil on Find
	articleRepositoryMock.On("GetByTitle", item.Title).Return(nil)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  item.Title,
		Source: f.Title,
		Author: "Unknown",
		Link:   item.Link,
		Tags:   item.Categories,
	}

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		feedId,
		item.Title,
		f.Link,
		item.Authors[0].Name,
		item.Link,
		item.Categories,
	).Return(&articleExpected)

	err := p.execute(articleRepositoryMock, feedId)

	articleRepositoryMock.AssertCalled(t, "GetByTitle", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", feedId, item.Title, f.Link, item.Authors[0].Name, item.Link, item.Categories)

	assert.Nil(t, err, "The error object should be nil")
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
	articleRepositoryMock := new(repository.MocksArticleRepository)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  item.Title,
		Source: f.Title,
		Author: "Unknown",
		Link:   item.Link,
		Tags:   []string{"article"},
	}

	// Article repository return article on Find
	articleRepositoryMock.On("GetByTitle", item.Title).Return(&articleExpected)

	err := p.execute(articleRepositoryMock, feedId)

	articleRepositoryMock.AssertCalled(t, "GetByTitle", item.Title)
	articleRepositoryMock.AssertNotCalled(t, "Create")

	assert.Nil(t, err, "The error object should be nil")
}

func TestExecuteNoCategoriesInItem(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{},
	}

	item := &gofeed.Item{
		Title:      "Article",
		Published:  "2023-06-04",
		Link:       "https://test.com/article-1",
		Categories: []string{},
		Authors: []*gofeed.Person{
			{
				Name: "Victor",
			},
		},
	}

	f := &gofeed.Feed{
		Title: "Super Test",
		Link:  "https://feed.link",
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
	articleRepositoryMock := new(repository.MocksArticleRepository)

	// Article repository return nil on Find
	articleRepositoryMock.On("GetByTitle", item.Title).Return(nil)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  item.Title,
		Source: f.Title,
		Author: "Unknown",
		Link:   item.Link,
		Tags:   []string{"article"},
	}

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		feedId,
		item.Title,
		f.Link,
		item.Authors[0].Name,
		item.Link,
		[]string{"article"},
	).Return(&articleExpected)

	err := p.execute(articleRepositoryMock, feedId)

	articleRepositoryMock.AssertCalled(t, "GetByTitle", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", feedId, item.Title, f.Link, item.Authors[0].Name, item.Link, []string{"article"})

	assert.Nil(t, err, "The error object should be nil")
}

func TestExecuteUnknownAuthor(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{},
	}

	item := &gofeed.Item{
		Title:      "Super Title for an Article",
		Published:  "2023-06-04",
		Link:       "https://test.com/article-1",
		Categories: []string{"tag2"},
		Authors:    []*gofeed.Person{},
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
	articleRepositoryMock := new(repository.MocksArticleRepository)

	// Article repository return nil on Find
	articleRepositoryMock.On("GetByTitle", item.Title).Return(nil)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  item.Title,
		Source: f.Title,
		Author: "Unknown",
		Link:   item.Link,
		Tags:   item.Categories,
	}

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		feedId,
		item.Title,
		f.Link,
		"Unknown",
		item.Link,
		item.Categories,
	).Return(&articleExpected)

	err := p.execute(articleRepositoryMock, feedId)

	articleRepositoryMock.AssertCalled(t, "GetByTitle", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", feedId, item.Title, f.Link, "Unknown", item.Link, item.Categories)

	assert.Nil(t, err, "The error object should be nil")
}

func TestExecuteWithCategorySpace(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{},
	}

	item := &gofeed.Item{
		Title:      "Super Title for an Article",
		Published:  "2023-06-04",
		Link:       "https://test.com/article-1",
		Categories: []string{"tag2 btc"},
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
	articleRepositoryMock := new(repository.MocksArticleRepository)

	// Article repository return nil on Find
	articleRepositoryMock.On("GetByTitle", item.Title).Return(nil)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  item.Title,
		Source: f.Title,
		Author: item.Authors[0].Name,
		Link:   item.Link,
		Tags:   []string{"tag2", "btc"},
	}

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		feedId,
		item.Title,
		f.Link,
		item.Authors[0].Name,
		item.Link,
		[]string{"tag2", "btc"},
	).Return(&articleExpected)

	err := p.execute(articleRepositoryMock, feedId)

	articleRepositoryMock.AssertCalled(t, "GetByTitle", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", feedId, item.Title, f.Link, item.Authors[0].Name, item.Link, []string{"tag2", "btc"})

	assert.Nil(t, err, "The error object should be nil")
}

func TestExecuteWithCategoryComa(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{},
	}

	item := &gofeed.Item{
		Title:      "Super Title for an Article",
		Published:  "2023-06-04",
		Link:       "https://test.com/article-1",
		Categories: []string{"tag2,btc"},
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
	articleRepositoryMock := new(repository.MocksArticleRepository)

	// Article repository return nil on Find
	articleRepositoryMock.On("GetByTitle", item.Title).Return(nil)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  item.Title,
		Source: f.Title,
		Author: item.Authors[0].Name,
		Link:   item.Link,
		Tags:   []string{"tag2", "btc"},
	}

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		feedId,
		item.Title,
		f.Link,
		item.Authors[0].Name,
		item.Link,
		[]string{"tag2", "btc"},
	).Return(&articleExpected)

	err := p.execute(articleRepositoryMock, feedId)

	articleRepositoryMock.AssertCalled(t, "GetByTitle", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", feedId, item.Title, f.Link, item.Authors[0].Name, item.Link, []string{"tag2", "btc"})

	assert.Nil(t, err, "The error object should be nil")
}

func TestExecuteWithUpperCaseCategory(t *testing.T) {
	p := &parser{
		url:  "https://test.com/feed",
		tags: []string{},
	}

	item := &gofeed.Item{
		Title:      "Super Title for an Article",
		Published:  "2023-06-04",
		Link:       "https://test.com/article-1",
		Categories: []string{"tag2", "BTC"},
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
	articleRepositoryMock := new(repository.MocksArticleRepository)

	// Article repository return nil on Find
	articleRepositoryMock.On("GetByTitle", item.Title).Return(nil)

	feedId := uuid.New().String()
	articleExpected := repository.DBArticle{
		Id:     uuid.New().String(),
		FeedId: feedId,
		Title:  item.Title,
		Source: f.Title,
		Author: item.Authors[0].Name,
		Link:   item.Link,
		Tags:   []string{"tag2", "btc"},
	}

	// Article repository return an article on Create
	articleRepositoryMock.On("Create",
		feedId,
		item.Title,
		f.Link,
		item.Authors[0].Name,
		item.Link,
		[]string{"tag2", "btc"},
	).Return(&articleExpected)

	err := p.execute(articleRepositoryMock, feedId)

	articleRepositoryMock.AssertCalled(t, "GetByTitle", item.Title)
	articleRepositoryMock.AssertCalled(t, "Create", feedId, item.Title, f.Link, item.Authors[0].Name, item.Link, []string{"tag2", "btc"})

	assert.Nil(t, err, "The error object should be nil")
}

func TestCleanCategoriesWithTitleContainsWeirdText(t *testing.T) {
	res := cleanCategories([]string{
		"The little test and green pepper",
	})

	assert.Equal(t, []string{"little", "test", "green", "pepper"}, res)
}
