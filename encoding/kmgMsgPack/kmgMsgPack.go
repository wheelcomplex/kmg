package kmgMsgPack

import (
	"bytes"
	"github.com/ugorji/go/codec"
	"io/ioutil"
	"os"
)

func WriteFile(path string, obj interface{}) (err error) {
	mh := codec.MsgpackHandle{}
	mh.AsSymbols = codec.AsSymbolNone
	mh.RawToString = true
	buf := &bytes.Buffer{}
	encoder := codec.NewEncoder(buf, &mh)
	err = encoder.Encode(obj)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, buf.Bytes(), os.FileMode(0777))
	return
}
