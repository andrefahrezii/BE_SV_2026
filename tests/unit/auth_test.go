package service_test

import (
	"testing"

	"github.com/you/sharing-vision-backend-v2/internal/auth"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	hash, err := auth.HashPassword("mypassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "mypassword", hash)
}

func TestComparePassword_Correct(t *testing.T) {
	hash, _ := auth.HashPassword("mypassword")
	assert.True(t, auth.ComparePassword(hash, "mypassword"))
}

func TestComparePassword_Incorrect(t *testing.T) {
	hash, _ := auth.HashPassword("mypassword")
	assert.False(t, auth.ComparePassword(hash, "wrongpassword"))
}

func TestGenerateToken(t *testing.T) {
	token, err := auth.GenerateToken(1, "test@example.com", "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
