package kmgGob

import (
	"bytes"
	"encoding/gob"
	"os"
)

func WriteFile(path string, obj interface{}) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		return
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	return encoder.Encode(obj)
}

func ReadFile(path string, obj interface{}) (err error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(0666))
	if err != nil {
		return
	}
	defer f.Close()
	encoder := gob.NewDecoder(f)
	return encoder.Decode(obj)
}

func Marshal(obj interface{}) (out []byte, err error) {
	b := &bytes.Buffer{}
	encoder := gob.NewEncoder(b)
	err = encoder.Encode(obj)
	if err != nil {
		return
	}
	return b.Bytes(), nil
}

func Unmarshal(data []byte, obj interface{}) (err error) {
	b := bytes.NewBuffer(data)
	encoder := gob.NewDecoder(b)
	err = encoder.Decode(obj)
	return
}
