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
	session Session
}

func (this *TestHttpHandlerService) TestFunc1(
	apiInput *struct{ A int },
	apiOutput *struct{ B int },
) error {
	apiOutput.B = apiInput.A + 1
	this.session.GetStore().Set("A",apiInput.A)
	return nil
}
func (this *TestHttpHandlerService) TestFunc1(
	apiInput *struct{ C int },
	apiOutput *struct{ D int },
) error {
	apiOutput.D = apiInput.C + this.session.GetStore().Get("A").(int)
	return nil
}

func TestHttpHandler(ot *testing.T) {
	t := test.NewTestTools(ot)
	c := dependencyInjection.NewContainer()
	err := c.Set("TestService", &TestHttpHandlerService{}, "")
	t.Equal(err, nil)
	apiManager := NewApiManagerFromContainer(c)
	h := &JsonHttpHandler{ApiManager: apiManager}

	output:=apiCall(h,t,`{"Name":"TestService.TestFunc1","Data":{"A":5}}`)
	t.Equal(output["Err"].(string), "")
	t.Equal(output["Data"].(map[string]interface{})["B"].(float64), 6)

	output:=apiCall(h,t,`{"Name":"TestService.TestFunc2","Data":{"C":5}}`)
	t.Equal(output["Err"].(string), "")
	t.Equal(output["Data"].(map[string]interface{})["D"].(float64), 10)

}

func apiCall(h *JsonHttpHandler,t *test.TestTools,json string)map[string]interface{}{
	w := httptest.NewRecorder()
	request, err := http.NewRequest("POST",
		"http://example.com/",
		bytes.NewBufferString(json))
	t.Equal(err, nil)

	h.ServeHTTP(w, request)
	var outputi interface{}
	err = json.Unmarshal(w.Body.Bytes(), &outputi)
	t.Equal(err, nil)
	return outputi.(map[string]interface{})
}
