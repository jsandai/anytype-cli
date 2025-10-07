package core

import (
	"encoding/base64"
	"strings"

	"github.com/anyproto/anytype-cli/core/output"
	"github.com/zalando/go-keyring"
)

const (
	keyringService        = "anytype-cli"
	keyringTokenUser      = "session-token"
	keyringBotAccountUser = "bot-account-key"
)

func SaveToken(token string) error {
	return keyring.Set(keyringService, keyringTokenUser, token)
}

func GetStoredToken() (string, error) {
	token, err := keyring.Get(keyringService, keyringTokenUser)
	if err != nil {
		return "", err
	}
	output.Info("Retrieved token from keyring %q", token)
	// Handle go-keyring base64 encoding that may be added on some platforms
	if strings.HasPrefix(token, "go-keyring-base64:") {
		encoded := strings.TrimPrefix(token, "go-keyring-base64:")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", err
		}
		token = string(decoded)
	}
	return token, nil
}

func DeleteStoredToken() error {
	return keyring.Delete(keyringService, keyringTokenUser)
}

func SaveBotAccountKey(accountKey string) error {
	return keyring.Set(keyringService, keyringBotAccountUser, accountKey)
}

func GetStoredBotAccountKey() (string, error) {
	key, err := keyring.Get(keyringService, keyringBotAccountUser)
	if err != nil {
		return "", err
	}
	output.Info("Retrieved key from keyring %q", key)
	// Handle go-keyring base64 encoding that may be added on some platforms
	if strings.HasPrefix(key, "go-keyring-base64:") {
		encoded := strings.TrimPrefix(key, "go-keyring-base64:")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", err
		}
		key = string(decoded)
	}
	return key, nil
}

func DeleteStoredBotAccountKey() error {
	return keyring.Delete(keyringService, keyringBotAccountUser)
}
