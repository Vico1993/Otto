package repository

import "github.com/stretchr/testify/mock"

type MocksArticleRepository struct {
	mock.Mock
}

func (m *MocksArticleRepository) Create(feedId string, title string, source string, author string, link string, tags []string) *DBArticle {
	args := m.Called(feedId, title, source, author, link, tags)
	return args.Get(0).(*DBArticle)
}

func (m *MocksArticleRepository) Delete(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}

func (m *MocksArticleRepository) GetByTitle(title string) *DBArticle {
	args := m.Called(title)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBArticle)
}

func (m *MocksArticleRepository) GetByFeedId(uuid string) []*DBArticle {
	args := m.Called(uuid)

	return args.Get(0).([]*DBArticle)
}

func (m *MocksArticleRepository) GetOne(uuid string) *DBArticle {
	args := m.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBArticle)
}

func (m *MocksArticleRepository) GetAll() []*DBArticle {
	args := m.Called()
	return args.Get(0).([]*DBArticle)
}

func (m *MocksArticleRepository) GetByChatAndTime(chatId string) []*DBArticle {
	args := m.Called(chatId)

	return args.Get(0).([]*DBArticle)
}
