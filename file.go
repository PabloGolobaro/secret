package secret

import (
	"encoding/json"
	"errors"
	"github.com/pablogolobaro/secret/encrypt"
	"io"
	"os"
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

func (f *FileVault) load() error {
	file, err := os.Open(f.filepath)
	if err != nil {
		f.keyValues = make(map[string]string)
		return nil
	}
	defer file.Close()

	r, err := encrypt.DecryptReader(f.encodingKey, file)
	if err != nil {
		return err
	}

	return f.readKeyValues(r)
}

func (f *FileVault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&f.keyValues)
}

func (f *FileVault) save() error {
	file, err := os.OpenFile(f.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil
	}
	defer file.Close()

	w, err := encrypt.EncryptWriter(f.encodingKey, file)
	if err != nil {
		return err
	}

	return f.writeKeyValues(w)
}

func (f *FileVault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(f.keyValues)
}

func (f *FileVault) Get(key string) (string, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	err := f.load()
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

	err := f.load()
	if err != nil {
		return err
	}

	f.keyValues[key] = value

	return f.save()
}
