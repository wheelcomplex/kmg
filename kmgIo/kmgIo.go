package kmgIo

import (
	"io"
	"io/ioutil"
)

func MustReadAll(r io.Reader) (b []byte) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}
