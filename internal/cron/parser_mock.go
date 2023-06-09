package cron

import (
	"github.com/Vico1993/Otto/internal/repository"
	"github.com/stretchr/testify/mock"
)

type parserMock struct {
	mock.Mock
}

func (m *parserMock) execute(articleRepository repository.IArticleRepository) (*parseResult, error) {
	args := m.Called(articleRepository)

	if args.Get(0) == nil {
		return nil, nil
	}

	if args.Error(0) != nil {
		return nil, args.Error(0)
	}

	return args.Get(0).(*parseResult), nil
}

func (m *parserMock) isCategoriesAndTagsMatch(categories []string) []string {
	args := m.Called(categories)

	return args.Get(0).([]string)
}
