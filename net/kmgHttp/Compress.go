package kmgHttp

import (
	"compress/flate"
	"compress/gzip"
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
		//w.Header().Set("Content-Encoding", "deflate")
		gzw, err := flate.NewWriter(w, -1)
		if err != nil {
			panic(err)
		}
		defer gzw.Close()
		gzr := ResponseWriterWraper{Writer: gzw, ResponseWriter: w}
		fn.ServeHTTP(gzr, r)
	})
}

func HttpHandleCompressGzipWrap(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			oldBody := r.Body
			defer oldBody.Close()
			var err error
			r.Body,err = gzip.NewReader(oldBody)
			if err!=nil{
				panic(err)
			}
		*/
		//w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		gzr := ResponseWriterWraper{Writer: gzw, ResponseWriter: w}
		fn.ServeHTTP(gzr, r)
	})
}
