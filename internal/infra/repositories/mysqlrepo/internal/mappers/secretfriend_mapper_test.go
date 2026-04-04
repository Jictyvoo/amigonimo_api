package mappers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

func TestStatusToEntity(t *testing.T) {
	cases := []struct {
		input    string
		expected entities.SecretFriendStatus
	}{
		{"draft", entities.StatusDraft},
		{"open", entities.StatusOpen},
		{"drawn", entities.StatusDrawn},
		{"closed", entities.StatusClosed},
		{"unknown_status", entities.SecretFriendStatus("unknown_status")},
		{"", entities.SecretFriendStatus("")},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := mappers.StatusToEntity(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}
