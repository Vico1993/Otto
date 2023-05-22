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
