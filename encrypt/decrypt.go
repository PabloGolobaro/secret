package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"io"
)

func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBDecrypter(block, iv), nil
}

func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to read the full iv")
	}

	stream, err := decryptStream(key, iv)
	if err != nil {
		return nil, err
	}

	return &cipher.StreamReader{S: stream, R: r}, nil
}

func Decrypt(key, cipherHex string) (string, error) {
	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("cipher too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream, err := decryptStream(key, iv)
	if err != nil {
		return "", err
	}

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil

}
