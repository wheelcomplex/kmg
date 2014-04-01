package kmgGob

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

func UrlBase64Marshal(obj interface{}) (out []byte, err error) {
	out1, err := Marshal(obj)
	if err != nil {
		return
	}
	out = make([]byte, base64.URLEncoding.EncodedLen(len(out1)))
	base64.URLEncoding.Encode(out, out1)
	return
}

func UrlBase64Unmarshal(data []byte, obj interface{}) (err error) {
	data1 := make([]byte, base64.URLEncoding.DecodedLen(len(data)))
	_, err = base64.URLEncoding.Decode(data1, data)
	if err != nil {
		return
	}
	b := bytes.NewBuffer(data1)
	encoder := gob.NewDecoder(b)
	err = encoder.Decode(obj)
	return
}
