package kmgGob

import (
	"bytes"
	"encoding/gob"
	"os"
)

func WriteFile(path string, obj interface{}) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0777))
	if err != nil {
		return
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	return encoder.Encode(obj)
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
