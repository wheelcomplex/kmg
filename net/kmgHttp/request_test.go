package kmgHttp

import (
	"github.com/bronze1man/kmg/kmgTest"
	"net/http"
	"net/url"
	"testing"
)

func TestGetStringFromGet(T *testing.T) {
	if "b" != getMockRequest().GetStringFromGet("a") {
		T.Fatalf("TestGetStringFromGet fail")
	}
	if "" != getMockRequest().GetStringFromGet("b") {
		T.Fatalf("TestGetStringFromGet fail")
	}
	if "" != getMockRequest().GetStringFromGet("") {
		T.Fatalf("TestGetStringFromGet fail")
	}
}

func TestGetGetStringMap(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	t.Equal(map[string]string{"a": "b", "c": "e"}, getMockRequest().GetGetStringMap())
	t.Equal(map[string]string{}, getEmptyQueryRequest().GetGetStringMap())
}

func getMockRequest() *KmgRequest {
	return &KmgRequest{&http.Request{
		URL: &url.URL{
			RawQuery: url.Values{
				"a": []string{"b", "c"},
				"c": []string{"e"},
			}.Encode()}},
	}
}

func getEmptyQueryRequest() *KmgRequest {
	return &KmgRequest{&http.Request{
		URL: &url.URL{
			RawQuery: url.Values{}.Encode()}},
	}
}
