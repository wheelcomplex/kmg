package webTypeAdmin

import (
	"encoding/json"
	"fmt"
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"net/http"
	//"reflect"
	"strings"
)

type Manager struct {
	context
}

func NewManager(ptr interface{}) (manager *Manager, err error) {
	ctx, err := newContext(ptr)
	if err != nil {
		return
	}
	manager = &Manager{context: *ctx}
	return
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//var err error
	pathS := req.FormValue("p")
	path := kmgType.ParsePath(pathS)
	switch req.Method {
	case "GET":
		s, err := manager.page(path)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(s))
	case "POST":
		f := strings.ToLower(req.FormValue("f"))
		var err error
		switch f {
		case "save":
			value := req.FormValue("v")
			err = manager.SaveByPath(path, value)
		case "delete":
			err = manager.DeleteByPath(path)
		default:
			err = fmt.Errorf("not support request function %s", f)
		}
		var errS string
		if err != nil {
			errS = err.Error()
		}
		out, err := json.Marshal(struct{ Err string }{Err: errS})
		if err != nil {
			panic(err)
		}
		w.Write(out)
		return
	default:
		w.Write([]byte(fmt.Sprintf("not support request method %s", req.Method)))
	}

	return
}

//show a page on some path
func (manager *Manager) page(path kmgType.Path) (html template.HTML, err error) {
	v, t, err := manager.GetElemByPath(path)
	if err != nil {
		return
	}
	html, err = t.HtmlView(v)
	if err != nil {
		return
	}
	return theTemplate.ExecuteNameToHtml("Main", struct {
		Path string
		Html template.HTML
	}{
		Path: path.String(),
		Html: html,
	})
}
