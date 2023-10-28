package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
)

type Encrypt struct {
	Nonce []byte
	AEAD  cipher.AEAD
}

func Get() (Encrypt, error) {
	password := "SECRET_PASSWORD"
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("aes.NewCipher error: %v\n", err)
		return Encrypt{}, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("cipher.NewGCM error: %v\n", err)
		return Encrypt{}, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	return Encrypt{
		Nonce: nonce,
		AEAD:  aesgcm,
	}, nil
}
