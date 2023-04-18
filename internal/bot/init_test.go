package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitSuccesfull(t *testing.T) {
	assert.Empty(t, ListCmd, "ListCmd should be empty at start")

	initCommand()

	assert.NotEmpty(t, ListCmd, "ListCmd should not be empty once initiated")
}
