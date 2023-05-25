package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMessage(t *testing.T) {
	// Set up test data
	title := "Title"
	platform := "platform"
	author := "author"
	tags := []string{"tag1", "tag2"}
	link := "https://example.com"

	// Override templates
	templates = []string{"$TITLE$ on $PLATFORM$. by $AUTHOR$. check link: $LINK$ is about: $TAGS$"}

	res := BuildMessage(title, platform, author, tags, link)
	assert.Equal(
		t,
		"Title on platform. by author. check link: https://example.com is about: #tag1, #tag2",
		res,
		"The template is not matching the expected translation!",
	)
}
