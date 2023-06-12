package repository

import (
	"github.com/Vico1993/Otto/internal/database"
	"github.com/stretchr/testify/mock"
)

type MocksArticleRep struct {
	mock.Mock
}

func (m *MocksArticleRep) Create(title string, published string, link string, source string, author string, match []string, tags ...string) *database.Article {
	args := m.Called(title, published, link, source, author, match, tags)
	return args.Get(0).(*database.Article)
}

func (m *MocksArticleRep) Find(key string, val string) *database.Article {
	args := m.Called(key, val)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*database.Article)
}
