package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hemagome/hefesto/config"
)

const tokenFile = "tokens.json"

var (
	encryptionKey []byte
	keyOnce       sync.Once
)

func initEncryptionKey() {
	keyOnce.Do(func() {
		cfg := config.GetConfig()
		if cfg == nil {
			panic("configuration not loaded")
		}

		key := cfg.Encryption.Key
		if len(key) != 32 {
			panic("encryption key must be exactly 32 bytes long")
		}

		encryptionKey = []byte(key)
	})
}

// encrypt cifra un texto con AES-256 GCM
func encrypt(data string) (string, error) {
	initEncryptionKey()

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encrypted := aesGCM.Seal(nonce, nonce, []byte(data), nil)
	return fmt.Sprintf("%x", encrypted), nil
}

// decrypt descifra un texto en AES-256 GCM
func decrypt(encryptedData string) (string, error) {
	initEncryptionKey()

	data, err := decodeHex(encryptedData)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	nonce, ciphertext := data[:12], data[12:]
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// decodeHex convierte una cadena hexadecimal en bytes
func decodeHex(str string) ([]byte, error) {
	bytes := make([]byte, len(str)/2)
	_, err := fmt.Sscanf(str, "%x", &bytes)
	return bytes, err
}

// StoreToken almacena un token cifrado en el archivo JSON
func StoreToken(provider, token string) error {
	encryptedToken, err := encrypt(token)
	if err != nil {
		return err
	}

	tokens, err := loadTokens()
	if err != nil {
		return err
	}

	tokens[provider] = encryptedToken
	return saveTokens(tokens)
}

// GetToken recupera un token descifrado
func GetToken(provider string) (string, error) {
	tokens, err := loadTokens()
	if err != nil {
		return "", err
	}

	encryptedToken, exists := tokens[provider]
	if !exists {
		return "", errors.New("token no encontrado")
	}

	return decrypt(encryptedToken)
}

// loadTokens carga los tokens desde el archivo
func loadTokens() (map[string]string, error) {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	var tokens map[string]string
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// saveTokens guarda los tokens en el archivo
func saveTokens(tokens map[string]string) error {
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tokenFile, data, 0600)
}
