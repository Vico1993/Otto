package database

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestTransformUUIDToString(t *testing.T) {
	res := TransformUUIDToString(
		pgtype.UUID{
			Bytes: [16]byte{149, 253, 186, 248, 6, 27, 79, 45, 156, 202, 148, 114, 19, 86, 128, 4},
			Valid: true,
		},
	)

	assert.Equal(t, "95fdbaf8-061b-4f2d-9cca-947213568004", res, "Should transform into the pgtype.UUID into a valid string")
}
