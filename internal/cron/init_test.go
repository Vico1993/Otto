package cron

import (
	"fmt"
	"testing"

	v2 "github.com/Vico1993/Otto/internal/repository/v2"
	"github.com/stretchr/testify/assert"
)

func TestDelayWith1Element(t *testing.T) {
	res := getDelay(1)

	assert.Equal(t, 60, res, "If only one element is passed a delay of 60 should be returned")
}

func TestDelayWith20Element(t *testing.T) {
	res := getDelay(20)

	assert.Equal(t, 3, res, "If 20 element are sent, should return 3")
}

func TestInit(t *testing.T) {
	jobs := Scheduler.Jobs()
	assert.Len(t, jobs, 0, "Should have 0 job in the queue at start")

	Init()

	jobs = Scheduler.Jobs()
	tags := Scheduler.GetAllTags()
	assert.Len(t, tags, 1, "Should have 1 tags from jobs in the queue at the end")
	assert.Len(t, jobs, 1, "Should have 1 job in the queue at the end")

	// Reset jobs
	_ = Scheduler.RemoveByTag(mainTag)
}

func TestCheckResetFeedNoFeedsInDB(t *testing.T) {
	feedRepositoryMock := new(v2.MocksFeedRepository)

	v2.Feed = feedRepositoryMock

	feedRepositoryMock.On("GetAll").Return([]*v2.DBFeed{})

	checkResetFeed()

	jobs, _ := Scheduler.FindJobsByTag(feedsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 0)
	feedRepositoryMock.AssertCalled(t, "GetAll")
	assert.Nil(t, jobs, "Should have no job link to "+feedsTag)
}

func TestCheckResetSameNumberOfFeed(t *testing.T) {
	_, _ = Scheduler.Every(1).Tag(feedsTag).Week().Do(func() {
		fmt.Println("Test job")
	})

	feedRepositoryMock := new(v2.MocksFeedRepository)

	v2.Feed = feedRepositoryMock

	feed := v2.DBFeed{
		Id:  "1234",
		Url: "https://google.com",
	}

	feedRepositoryMock.On("GetAll").Return([]*v2.DBFeed{&feed})

	jobs, _ := Scheduler.FindJobsByTag(feedsTag)
	assert.Len(t, jobs, 1, "Should have 1 job in the queue at start")

	checkResetFeed()

	jobs, _ = Scheduler.FindJobsByTag(feedsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 1)
	feedRepositoryMock.AssertCalled(t, "GetAll")
	assert.Len(t, jobs, 1, "Should have 1 job in the queue after call")

	// Reset jobs
	_ = Scheduler.RemoveByTag(feedsTag)
}

func TestCheckResetAddJob(t *testing.T) {
	_, _ = Scheduler.Every(1).Tag(feedsTag).Week().Do(func() {
		fmt.Println("Test job")
	})

	feedRepositoryMock := new(v2.MocksFeedRepository)

	v2.Feed = feedRepositoryMock

	feed1 := v2.DBFeed{
		Id:  "1234",
		Url: "https://google1.com",
	}
	feed2 := v2.DBFeed{
		Id:  "5678",
		Url: "https://google2.com",
	}

	feedRepositoryMock.On("GetAll").Return([]*v2.DBFeed{&feed1, &feed2})

	jobs, _ := Scheduler.FindJobsByTag(feedsTag)
	assert.Len(t, jobs, 1, "Should have 1 job in the queue at start")

	checkResetFeed()

	jobs, _ = Scheduler.FindJobsByTag(feedsTag)
	tags := Scheduler.GetAllTags()

	assert.Len(t, tags, 2)
	assert.Equal(t, []string{"feed", "feed"}, tags, "Should have 2 tags call feeds")
	feedRepositoryMock.AssertCalled(t, "GetAll")
	assert.Len(t, jobs, 2, "Should have 2 job in the queue after call")

	// Reset jobs
	_ = Scheduler.RemoveByTag(feedsTag)
}
