package cron

import (
	"fmt"
	"testing"

	"github.com/Vico1993/Otto/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCheckResetFeedNoFeedsInDB(t *testing.T) {
	feedRepositoryMock := new(repository.MocksFeedRepository)

	repository.Feed = feedRepositoryMock

	feedRepositoryMock.On("GetAllActive").Return([]*repository.DBFeed{})

	checkFeed()

	jobs, _ := Scheduler.FindJobsByTag(feedsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 0)
	feedRepositoryMock.AssertCalled(t, "GetAllActive")
	assert.Nil(t, jobs, "Should have no job link to "+feedsTag)
}

func TestCheckResetSameNumberOfFeed(t *testing.T) {
	_, _ = Scheduler.Every(1).Tag(feedsTag).Week().Do(func() {
		fmt.Println("Test job")
	})

	feedRepositoryMock := new(repository.MocksFeedRepository)

	repository.Feed = feedRepositoryMock

	feed := repository.DBFeed{
		Id:       "1234",
		Url:      "https://google.com",
		Disabled: false,
	}

	feedRepositoryMock.On("GetAllActive").Return([]*repository.DBFeed{&feed})

	jobs, _ := Scheduler.FindJobsByTag(feedsTag)
	assert.Len(t, jobs, 1, "Should have 1 job in the queue at start")

	checkFeed()

	jobs, _ = Scheduler.FindJobsByTag(feedsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 1)
	feedRepositoryMock.AssertCalled(t, "GetAllActive")
	assert.Len(t, jobs, 1, "Should have 1 job in the queue after call")

	// Reset jobs
	_ = Scheduler.RemoveByTag(feedsTag)
}

func TestCheckResetAddJob(t *testing.T) {
	_, _ = Scheduler.Every(1).Tag(feedsTag).Week().Do(func() {
		fmt.Println("Test job")
	})

	feedRepositoryMock := new(repository.MocksFeedRepository)

	repository.Feed = feedRepositoryMock

	feed1 := repository.DBFeed{
		Id:       "1234",
		Url:      "https://google1.com",
		Disabled: false,
	}
	feed2 := repository.DBFeed{
		Id:       "5678",
		Url:      "https://google2.com",
		Disabled: false,
	}

	feedRepositoryMock.On("GetAllActive").Return([]*repository.DBFeed{&feed1, &feed2})

	jobs, _ := Scheduler.FindJobsByTag(feedsTag)
	assert.Len(t, jobs, 1, "Should have 1 job in the queue at start")

	checkFeed()

	jobs, _ = Scheduler.FindJobsByTag(feedsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 2)
	assert.Equal(t, []string{"feed", "feed"}, tags, "Should have 2 tags call feeds")
	feedRepositoryMock.AssertCalled(t, "GetAllActive")
	assert.Len(t, jobs, 2, "Should have 2 job in the queue after call")

	// Reset jobs
	_ = Scheduler.RemoveByTag(feedsTag)
}
