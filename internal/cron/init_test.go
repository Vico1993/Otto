package cron

import (
	"testing"

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
