package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenGeneration(t *testing.T) {
	accessToken, refreshToken, err := GenerateTokens("some-user-id", "127.0.0.1")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}
