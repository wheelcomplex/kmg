package ajkApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bronze1man/kmg/dependencyInjection"
	"github.com/bronze1man/kmg/sessionStore"
	"github.com/bronze1man/kmg/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

var globalTestVar = 0

type TestHttpHandlerService struct {
	session *Session
}

func (this *TestHttpHandlerService) TestFunc1(
	apiInput *struct{ A int },
	apiOutput *struct{ B int },
) error {
	apiOutput.B = apiInput.A + 1
	store, err := this.session.GetStore()
	if err != nil {
		return err
	}
	store.Set("A", apiInput.A)
	return nil
}
func (this *TestHttpHandlerService) TestFunc2(
	apiInput *struct{ C int },
	apiOutput *struct{ D int },
) error {
	store, err := this.session.GetStore()
	if err != nil {
		return err
	}
	a, ok := store.Get("A")
	if !ok {
		return errors.New("A not exist")
	}
	apiOutput.D = apiInput.C + a.(int)
	return nil
}
func (this *TestHttpHandlerService) TestFunc3() {
	globalTestVar = 1
	return
}

func TestHttpHandler(ot *testing.T) {
	t := test.NewTestTools(ot)
	c := dependencyInjection.NewContainer()
	err := c.SetFactory("TestService", func(c *dependencyInjection.Container) (interface{}, error) {
		session, err := c.Get("session")
		if err != nil {
			fmt.Println(1)
			return nil, err
		}
		return &TestHttpHandlerService{session: session.(*Session)}, nil
	}, dependencyInjection.ScopeRequest)
	t.Equal(err, nil)
	apiManager := NewApiManagerFromContainer(c)
	h := &JsonHttpHandler{ApiManager: apiManager,
		SessionStoreManager: &sessionStore.Manager{sessionStore.NewMemoryProvider()},
	}

	output := apiCall(h, t, `{"Name":"TestService.TestFunc1","Data":{"A":5}}`)
	t.Equal(output["Err"].(string), "")
	t.Equal(output["Data"].(map[string]interface{})["B"].(float64), 6.0)

	output = apiCall(h, t, `{"Name":"TestService.TestFunc2","Data":{"C":5},"Guid":"`+
		output["Guid"].(string)+`"}`)
	t.Equal(output["Err"].(string), "")
	t.Equal(output["Data"].(map[string]interface{})["D"].(float64), 10.0)

	globalTestVar = 2
	apiCall(h, t, `{"Name":"TestService.TestFunc3","Data":null,"Guid":""}`)
	t.Equal(globalTestVar, 1)

}

func apiCall(h *JsonHttpHandler, t *test.TestTools, j string) map[string]interface{} {
	w := httptest.NewRecorder()
	request, err := http.NewRequest("POST",
		"http://example.com/",
		bytes.NewBufferString(j))
	t.Equal(err, nil)

	h.ServeHTTP(w, request)
	var outputi interface{}
	err = json.Unmarshal(w.Body.Bytes(), &outputi)
	t.Equal(err, nil)
	return outputi.(map[string]interface{})
}
