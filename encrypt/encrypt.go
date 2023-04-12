package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	io.WriteString(hasher, key)
	cipherKey := hasher.Sum(nil)

	return aes.NewCipher(cipherKey)
}

func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBEncrypter(block, iv), nil
}

func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream, err := encryptStream(key, iv)
	if err != nil {
		return nil, err
	}

	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write full iv to writer")
	}

	return &cipher.StreamWriter{S: stream, W: w}, err
}

func Encrypt(key, plaintext string) (string, error) {
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream, err := encryptStream(key, iv)
	if err != nil {
		return "", err
	}

	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}
