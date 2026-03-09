package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// @func decrypt
// @description AES-256-GCM 암호문을 복호화한다

type DecryptInput struct {
	Ciphertext string // base64 인코딩
	Key        string // 32바이트 hex
}

type DecryptOutput struct {
	Plaintext string
}

func Decrypt(in DecryptInput) (DecryptOutput, error) {
	data, err := base64.StdEncoding.DecodeString(in.Ciphertext)
	if err != nil {
		return DecryptOutput{}, err
	}
	keyBytes, err := hex.DecodeString(in.Key)
	if err != nil {
		return DecryptOutput{}, err
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return DecryptOutput{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return DecryptOutput{}, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return DecryptOutput{}, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return DecryptOutput{}, err
	}
	return DecryptOutput{Plaintext: string(plaintext)}, nil
}
