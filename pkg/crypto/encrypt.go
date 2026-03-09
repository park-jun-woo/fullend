package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// @func encrypt
// @description 평문을 AES-256-GCM으로 암호화한다

type EncryptInput struct {
	Plaintext string
	Key       string // 32바이트 hex
}

type EncryptOutput struct {
	Ciphertext string // base64 인코딩
}

func Encrypt(in EncryptInput) (EncryptOutput, error) {
	keyBytes, err := hex.DecodeString(in.Key)
	if err != nil {
		return EncryptOutput{}, err
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return EncryptOutput{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptOutput{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return EncryptOutput{}, err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(in.Plaintext), nil)
	return EncryptOutput{Ciphertext: base64.StdEncoding.EncodeToString(sealed)}, nil
}
