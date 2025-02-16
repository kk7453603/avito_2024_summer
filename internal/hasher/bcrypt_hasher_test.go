package hasher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_Hash(t *testing.T) {
	h := New()

	passwd1 := "best-password-1"
	passwd2 := "best-password-2"

	hashedPasswd1, err1 := h.Hash(passwd1)
	require.NoError(t, err1, "Password #1 hashing failure")
	require.NotEmpty(t, hashedPasswd1, "Hashed password #1 must not be empty")

	hashedPasswd2, err2 := h.Hash(passwd2)
	require.NoError(t, err2, "Password #2 hashing failure")
	require.NotEmpty(t, hashedPasswd2, "Hashed password #2 must not be empty")

	require.NotEqual(t, passwd1, hashedPasswd1, "Hash #1 must not match the original password #1")
	require.NotEqual(t, passwd2, hashedPasswd2, "Hash #2 must not match the original password #2")

	require.NotEqual(t, hashedPasswd1, hashedPasswd2, "Hashes of different passwords must be unique")
}

func TestBcryptHasher_HashError(t *testing.T) {
	h := New()

	// len(passwd) > 72 | = 73
	passwd := "ZcFIBYnZqaYsSTGKpqUrIQdxKeHdlqZItjcrobEIsgkigmVHqyFkwDxuCrgfcBRUdDvwaYyxC"
	hashedPasswd, err := h.Hash(passwd)

	require.Error(t, err, "Password hashing no failure")
	require.Empty(t, hashedPasswd, "Hashed password should be empty")
}

func TestBcryptHasher_Compare(t *testing.T) {
	h := New()

	passwd := "best-passwd"
	hashedPasswd, _ := h.Hash(passwd)

	t.Run("Correct passwd", func(t *testing.T) {
		require.True(t, h.Compare(hashedPasswd, passwd), "Expected TRUE with the correct password")
	})

	t.Run("Incorrect passwd", func(t *testing.T) {
		require.False(t, h.Compare(hashedPasswd, "wrong-passwd"), "Expected FALSE with invalid password")
	})
}
