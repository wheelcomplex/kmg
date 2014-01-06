package kmgGob

import (
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
