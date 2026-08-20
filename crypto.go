package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// Lightweight AES-256-GCM encryption for provider credentials at rest.
// The key comes from PAYMENT_ENC_KEY (32 bytes). When unset, a fixed dev key
// is used so local testing works out of the box; production must set it.

var errInvalidCiphertext = errors.New("invalid ciphertext")

func encKey() []byte {
	s := os.Getenv("PAYMENT_ENC_KEY")
	if len(s) == 32 {
		return []byte(s)
	}
	// dev fallback — 32 bytes
	return []byte("payment-gateway-dev-key-0000000000000!!")[:32]
}

func encryptSecret(plain string) (string, error) {
	if plain == "" {
		return "", nil
	}
	block, err := aes.NewCipher(encKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(out), nil
}

func decryptSecret(enc string) (string, error) {
	if enc == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", errInvalidCiphertext
	}
	block, err := aes.NewCipher(encKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", errInvalidCiphertext
	}
	nonce, ct := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", errInvalidCiphertext
	}
	return string(plain), nil
}
