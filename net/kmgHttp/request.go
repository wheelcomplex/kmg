package kmgHttp

import "net/http"

var MaxMemory = 64 * 1024 * 1024

type Request struct {
	*http.Request
}

func NewRequestFromOrigin(req *http.Request) *Request {
	return &Request{Request: req}
}

//get first string from a key of url query ,if the key not exist return ""
// some thing like $_GET[$key] in php
func (r *Request) GetStringFromGet(key string) string {
	ret1 := r.Request.URL.Query()[key]
	if ret1 == nil {
		return ""
	}
	return ret1[0]
}

//return level 1 string map of url query
func (r *Request) GetGetStringMap() (output map[string]string) {
	query := r.Request.URL.Query()
	output = make(map[string]string)
	for k, v := range query {
		output[k] = v[0]
	}
	return
}
