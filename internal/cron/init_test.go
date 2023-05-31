package cron

import (
	"fmt"
	"testing"

	"github.com/Vico1993/Otto/internal/database"
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

func TestResetCronForChat(t *testing.T) {
	chat := database.NewChat("1234", 123, nil)
	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job")
	})

	assert.Len(t, scheduler.Jobs(), 1, "We should have only 1 job present at the initialisation")

	resetCronForChatId(chat)

	assert.Len(t, scheduler.Jobs(), 0, "After the delete we should have 0 jobs")
}

func TestAddJobForChat(t *testing.T) {
	feed := database.NewFeed("https://google.com")
	chat := database.NewChat("1234", 123, []database.Feed{
		*feed,
	})

	startJobForChat(chat)

	jobs := scheduler.Jobs()
	assert.Len(t, jobs, 1, "We should have only 1 job")

	job := jobs[0]
	assert.Len(t, job.Tags(), 1, "Only one tag should be attached to the job")
	assert.Equal(t, job.Tags()[0], chat.ChatId, "The first and only tag should be equal to our chatid")

	_ = scheduler.RemoveByTag(chat.ChatId)
}

func TestSetupCronForChat(t *testing.T) {
	feed := database.NewFeed("https://google.com")
	chat := database.NewChat("1234", 123, []database.Feed{
		*feed,
	})
	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job")
	})
	_, _ = scheduler.Every(1).Tag(chat.ChatId).Week().Do(func() {
		fmt.Println("Test job2")
	})

	fmt.Println(len(scheduler.Jobs()))

	assert.Len(t, scheduler.Jobs(), 2, "We should have only 2 job present at the initialisation")

	SetupCronForChat(chat)

	jobs := scheduler.Jobs()
	assert.Len(t, jobs, 1, "We should have only 1 job now for that chat")

	job := jobs[0]
	assert.Len(t, job.Tags(), 1, "Only one tag should be attached to the job")
	assert.Equal(t, job.Tags()[0], chat.ChatId, "The first and only tag should be equal to our chatid")

	_ = scheduler.RemoveByTag(chat.ChatId)
}
