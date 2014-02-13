package kmgHttp

import (
	"compress/flate"
	"io"
	"net/http"
)

type ResponseWriterWraper struct {
	io.Writer
	http.ResponseWriter
}

func (w ResponseWriterWraper) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// a flate(DEFLATE) compress wrap around http request and response,
// set http header
// !!not check http header!!
func HttpHandleCompressFlateWrap(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oldBody := r.Body
		defer oldBody.Close()
		r.Body = flate.NewReader(oldBody)
		w.Header().Set("Content-Encoding", "deflate")
		gzw, err := flate.NewWriter(w, -1)
		if err != nil {
			panic(err)
		}
		defer gzw.Close()
		gzr := ResponseWriterWraper{Writer: gzw, ResponseWriter: w}
		fn.ServeHTTP(gzr, r)
	})
}
