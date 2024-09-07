package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	t.Run("Valid Key Generation", func(t *testing.T) {
		key, err := GenerateKey()
		assert.NoError(t, err, "Expected no error during key generation")
		assert.Len(t, key, 64, "Expected key length of 64 characters")
	})

	t.Run("Randomness", func(t *testing.T) {
		// Generate multiple keys and verify they are different
		keys := make(map[string]struct{})
		for i := 0; i < 100; i++ {
			key, err := GenerateKey()
			assert.NoError(t, err, "Expected no error during key generation")
			_, exists := keys[key]
			assert.False(t, exists, "Expected unique keys")
			keys[key] = struct{}{}
		}
	})

	t.Run("Key Length", func(t *testing.T) {
		key, err := GenerateKey()
		assert.NoError(t, err, "Expected no error during key generation")

		decoded, err := hex.DecodeString(key)
		assert.NoError(t, err, "Expected no error during key decoding")
		assert.Len(t, decoded, 32, "Expected key length of 32 bytes")
	})
}
