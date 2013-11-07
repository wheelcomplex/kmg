package kmgHttp

import "net/http"
import "net/http/httputil"
import "bytes"
import "bufio"
import "net/url"

//sometimes it is hard to remember how to get response from bytes ...
func NewResponseFromBytes(r []byte) (resp *http.Response, err error) {
	return http.ReadResponse(bufio.NewReader(bytes.NewBuffer(r)), &http.Request{})
}

//sometimes it is hard to remember how to dump response to bytes
func DumpResponseToBytes(resp *http.Response) (b []byte, err error) {
	return httputil.DumpResponse(resp, true)
}

//sort url query to unique it
func SortUrlQuery(vurl *url.URL) {
	query := vurl.Query()
	vurl.RawQuery = query.Encode()
}
