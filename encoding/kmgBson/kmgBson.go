package kmgBson

import (
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"os"
)

func WriteFile(path string, obj interface{}) (err error) {
	out, err := bson.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, os.FileMode(0777))
}
