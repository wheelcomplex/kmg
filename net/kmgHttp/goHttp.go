package kmgHttp

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httputil"
)

func NewRequestFromByte(r []byte) (req *http.Request, err error) {
	return http.ReadRequest(bufio.NewReader(bytes.NewReader(r)))
}

//sometimes it is hard to remember how to get response from bytes ...
func NewResponseFromBytes(r []byte) (resp *http.Response, err error) {
	return http.ReadResponse(bufio.NewReader(bytes.NewBuffer(r)), &http.Request{})
}

//sometimes it is hard to remember how to dump response to bytes
func DumpResponseToBytes(resp *http.Response) (b []byte, err error) {
	return httputil.DumpResponse(resp, true)
}
