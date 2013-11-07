package ajkApi

import (
	"bytes"
	"kmg/dependencyInjection"
	"kmg/test"
	"net/http"
	"net/http/httptest"
	"testing"
	//"fmt"
	"encoding/json"
)

type TestHttpHandlerService struct {
}

func (this *TestHttpHandlerService) TestFunc1(apiInput *struct{ A int }, apiOutput *struct{ B int }) error {
	apiOutput.B = apiInput.A + 1
	return nil
}

func TestHttpHandler(ot *testing.T) {
	t := test.NewTestTools(ot)
	c := dependencyInjection.NewContainer()
	err := c.Set("TestService", &TestHttpHandlerService{}, "")
	t.Equal(err, nil)
	apiManager := NewApiManagerFromContainer(c)
	httpHandler := &JsonHttpHandler{ApiManager: apiManager}
	w := httptest.NewRecorder()
	request, err := http.NewRequest("POST",
		"http://example.com/",
		bytes.NewBufferString(`{"Name":"TestService.TestFunc1","Data":{"A":1}}`))
	t.Equal(err, nil)

	httpHandler.ServeHTTP(w, request)
	var outputi interface{}
	err = json.Unmarshal(w.Body.Bytes(), &outputi)
	t.Equal(err, nil)
	output := outputi.(map[string]interface{})
	t.Equal(output["Err"].(string), "")
	t.Equal(output["Data"].(map[string]interface{})["B"].(float64), 2.0)
}
