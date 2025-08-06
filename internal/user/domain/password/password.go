// Package password provides services for hashing and verifying passwords using Argon2.
package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// params holds the configuration for the Argon2 hashing algorithm.
// These parameters should be tuned based on the hardware and security requirements.
type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// DefaultParams provides sensible default parameters for Argon2id.
var DefaultParams = &params{
	memory:      64 * 1024, // 64 MB
	iterations:  1,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

// Hash creates a new Argon2id hash from a plaintext password.
// The salt and all parameters are encoded into the final hash string.
func Hash(password string) (string, error) {
	p := DefaultParams

	// Generate a cryptographically secure random salt.
	salt := make([]byte, p.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Generate the hash.
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Encode the salt and hash to Base64.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format the encoded hash string for storage.
	// This format is compatible with many other implementations.
	// Format: $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// Verify compares a plaintext password with a stored Argon2id hash.
// It returns true if the password and hash match, and false otherwise.
func Verify(password, encodedHash string) (bool, error) {
	// Parse the encoded hash string to extract parameters, salt, and the hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Generate a hash from the provided password using the exact same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Compare the two hashes in constant time to prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// decodeHash parses the encoded hash string and returns the parameters, salt, and hash.
func decodeHash(encodedHash string) (*params, []byte, []byte, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid encoded hash format")
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, fmt.Errorf("unsupported algorithm: %s", vals[1])
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid version: %w", err)
	}
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("incompatible version: got %d, want %d", version, argon2.Version)
	}

	p := &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid salt: %w", err)
	}
	p.saltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid hash: %w", err)
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
