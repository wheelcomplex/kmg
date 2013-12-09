package webTypeAdmin

import (
	"encoding/json"
	"fmt"
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"net/http"
	"reflect"
	"strings"
)

type Manager struct {
	ctx *context
}

func NewManagerFromPtr(ptr interface{}) (*Manager, error) {
	ctx := &context{}
	rt := reflect.TypeOf(ptr)
	if rt.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("NewManagerFromPtr need a ptr,but get %T", ptr)
	}
	ctx.rootValue = reflect.ValueOf(ptr)
	t, err := ctx.typeOfFromReflect(rt)
	if err != nil {
		return nil, err
	}
	ctx.rootType = t
	m := &Manager{ctx: ctx}
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
func (manager *Manager) page(path kmgType.Path) (string, error) {

	v, err := manager.getValueByPath(path)
	if err != nil {
		return "", err
	}
	t, err := manager.ctx.newTypeFromReflect(v.Type())
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
	return "", nil
}

func (manager *Manager) getValueByPath(p kmgType.Path) (v reflect.Value, err error) {
	t := manager.ctx.rootType
	v = manager.ctx.rootValue
	for _, ps := range p {
		v, err = t.getSubValueByString(v, ps)
		if err != nil {
			return reflect.Value{}, err
		}
		t, err = manager.ctx.newTypeFromReflect(v.Type())
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
	t, err := manager.ctx.newTypeFromReflect(v.Type())
	if err != nil {
		return err
	}
	return t.create(v, path[len(path)-1])
}

//set a value in some where in the whole type tree
//a workaround of MapIndex not addressable problem,
//if the reverse path come map first , it need reset data at map level
//if the reverse path come ptr first , it is ok
func (manager *Manager) save(path Path, value string) error {
	t := manager.ctx.rootType
	v := manager.ctx.rootValue
	var lastCanSetV reflect.Value
	var lastCanSetT typeInterface
	var lastCanSetI int
	var err error
	for i, ps := range path {
		v, err = t.getSubValueByString(v, ps)
		if err != nil {
			return err
		}
		t, err = manager.ctx.newTypeFromReflect(v.Type())
		if err != nil {
			return err
		}
		if v.CanSet() {
			lastCanSetV = v
			lastCanSetT = t
			lastCanSetI = i
		}
	}
	if lastCanSetV == v {
		return t.save(v, value)
	}
	if mt, ok := lastCanSetT.(*mapType); ok {
		return mt.mapSave(v, path[lastCanSetI+1:len(path)-1], value)
	}
	return fmt.Errorf("[manager.save] impossable code path")

	parentPath := path[:len(path)-1]
	v, err := manager.getValueByPath(parentPath)
	if err != nil {
		return err
	}
	t, err := manager.ctx.newTypeFromReflect(v.Type())
	if err != nil {
		return err
	}
	if mt, ok := t.(*mapType); ok {
		return mt.mapSave(v, path[len(path)-1], value)
	}
	v, err = t.getSubValueByString(v, path[len(path)-1])
	if err != nil {
		return err
	}
	t, err = manager.ctx.newTypeFromReflect(v.Type())
	if err != nil {
		return err
	}
	return t.save(v, value)

}

func (manager *Manager) save(path Path, value string) (err error) {
	pRootValue := &manager.ctx.rootValue
	err = manager.ctx.rootType.Save(pRootValue, path, value)
	if err != nil {
		return err
	}
	if pRootValue != &manager.ctx.rootValue {
		return fmt.Errorf("[manager.save] can not save")
	}
	return nil
}

//delete an object in some where in the whole type tree
//can delete a slice elem or map elem
func (manager *Manager) delete(path Path) (err error) {
	return nil

	parentPath := path[:len(path)-1]
	v, err := manager.getValueByPath(parentPath)
	if err != nil {
		return err
	}
	t, err := manager.ctx.newTypeFromReflect(v.Type())
	if err != nil {
		return err
	}
	return t.delete(v, path[len(path)-1])

}
