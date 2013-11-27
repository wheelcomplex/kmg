package webTypeAdmin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"
)

type Manager struct {
	rootType  typeInterface
	rootValue reflect.Value
}

func NewManagerFromPtr(ptr interface{}) (*Manager, error) {
	rt := reflect.TypeOf(ptr)
	if rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("NewManagerFromPtr need a ptr,but get %T", ptr)
	}
	t, err := newTypeFromReflect(rt)
	if err != nil {
		return nil, err
	}
	m := &Manager{rootType: t, rootValue: reflect.ValueOf(ptr)}
	return m, nil
}

func (manager *Manager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//var err error
	pathS := req.FormValue("p")
	path := parsePath(pathS)
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
		case "create":
			err = manager.create(path)
		case "save":
			value := req.FormValue("v")
			err = manager.save(path, value)
		case "delete":
			err = manager.delete(path)
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
func (manager *Manager) page(path Path) (string, error) {
	v, err := manager.getValueByPath(path)
	if err != nil {
		return "", err
	}
	t, err := newTypeFromReflect(v.Type())
	if err != nil {
		return "", err
	}
	b, err := theTemplate.ExecuteNameToByte("Main", struct {
		Path string
		Html template.HTML
	}{
		Path: path.String(),
		Html: t.Html(v),
	})
	return string(b), err
}

func (manager *Manager) getValueByPath(p Path) (v reflect.Value, err error) {
	t := manager.rootType
	v = manager.rootValue
	for _, ps := range p {
		v, err = t.getSubValueByString(v, ps)
		if err != nil {
			return reflect.Value{}, err
		}
		t, err = newTypeFromReflect(v.Type())
		if err != nil {
			return reflect.Value{}, err
		}
	}
	return v, nil
}

//create an empty object in some where in the whole type tree
//new key in path,or 0 if it is a slice
func (manager *Manager) create(path Path) error {
	parentPath := path[:len(path)-1]
	v, err := manager.getValueByPath(parentPath)
	if err != nil {
		return err
	}
	t, err := newTypeFromReflect(v.Type())
	if err != nil {
		return err
	}
	return t.create(v, path[len(path)-1])
}

//set a value in some where in the whole type tree
func (manager *Manager) save(path Path, value string) error {
	v, err := manager.getValueByPath(path)
	if err != nil {
		return err
	}
	t, err := newTypeFromReflect(v.Type())
	if err != nil {
		return err
	}
	return t.save(v, value)
}

//delete an object in some where in the whole type tree
//can delete a slice elem or map elem
func (manager *Manager) delete(path Path) (err error) {
	parentPath := path[:len(path)-1]
	v, err := manager.getValueByPath(parentPath)
	if err != nil {
		return err
	}
	t, err := newTypeFromReflect(v.Type())
	if err != nil {
		return err
	}
	return t.delete(v, path[len(path)-1])
}
