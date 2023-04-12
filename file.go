package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pablogolobaro/secret/encrypt"
	"io"
	"os"
	"strings"
	"sync"
)

func File(encodingKey, file string) *FileVault {
	return &FileVault{encodingKey: encodingKey, filepath: file, keyValues: make(map[string]string)}
}

type FileVault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

func (f *FileVault) loadKeyValues() error {
	file, err := os.Open(f.filepath)
	if err != nil {
		f.keyValues = make(map[string]string)
		return nil
	}
	defer file.Close()

	var sb strings.Builder

	_, err = io.Copy(&sb, file)
	if err != nil {
		return err
	}

	decryptedJSON, err := encrypt.Decrypt(f.encodingKey, sb.String())
	if err != nil {
		return err
	}

	r := strings.NewReader(decryptedJSON)

	dec := json.NewDecoder(r)

	err = dec.Decode(&f.keyValues)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileVault) saveKeyValues() error {
	var sb strings.Builder

	enc := json.NewEncoder(&sb)

	err := enc.Encode(f.keyValues)
	if err != nil {
		return err
	}

	encryptedJSON, err := encrypt.Encrypt(f.encodingKey, sb.String())
	if err != nil {
		return err
	}

	file, err := os.OpenFile(f.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, encryptedJSON)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileVault) Get(key string) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	err := f.loadKeyValues()
	if err != nil {
		return "", err
	}

	value, ok := f.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")

	}

	return value, nil
}

func (f *FileVault) Set(key, value string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	err := f.loadKeyValues()
	if err != nil {
		return err
	}

	f.keyValues[key] = value

	err = f.saveKeyValues()

	return nil
}
