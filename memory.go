package secret

import (
	"errors"
	"github.com/pablogolobaro/secret/encrypt"
)

func Memory(encodingKey string) MemoryVault {
	return MemoryVault{encodingKey: encodingKey, keyValues: make(map[string]string)}
}

type MemoryVault struct {
	encodingKey string
	keyValues   map[string]string
}

func (v *MemoryVault) Get(key string) (string, error) {
	hex, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")

	}
	decrypt, err := encrypt.Decrypt(v.encodingKey, hex)
	if err != nil {
		return "", err
	}
	return decrypt, nil
}

func (v *MemoryVault) Set(key, value string) error {
	encryptedValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}
	v.keyValues[key] = encryptedValue
	return nil
}
