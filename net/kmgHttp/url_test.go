package kmgHttp

import (
	"testing"
	"kmg/kmpTest"
)

func TestNewUrlByString(T *testing.T) {
	url,err:=NewUrlByString("http://www.google.com")
	kmpTest.Assert(T,nil,err)
	kmpTest.Assert(T,"http://www.google.com",url.String())

}
