package ajkApi

import (
	"compress/flate"
	"github.com/bronze1man/kmg/net/kmgHttp"
	"net/http"
)

type HttpApiFilter func(c *HttpApiContext, fc []HttpApiFilter)
type HttpApiContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	ApiName        string
}
type HttpApiFilterManager struct {
	Filters []HttpApiFilter
}

func (m *HttpApiFilterManager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := &HttpApiContext{
		Request:        req,
		ResponseWriter: w,
	}
	m.Filters[0](c, m.Filters[1:])
}

func HttpApiDeflateCompressFilter(c *HttpApiContext, fc []HttpApiFilter) {
	oldBody := c.Request.Body
	defer oldBody.Close()
	c.Request.Body = flate.NewReader(oldBody)
	gzw, err := flate.NewWriter(c.ResponseWriter, -1)
	if err != nil {
		panic(err)
	}
	defer gzw.Close()
	httpWriter := kmgHttp.ResponseWriterWraper{Writer: gzw, ResponseWriter: c.ResponseWriter}
	c.ResponseWriter = httpWriter
	fc[0](c, fc[1:])
}
