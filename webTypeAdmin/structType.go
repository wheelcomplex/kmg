package webTypeAdmin

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"html/template"
	"reflect"
)

//path -> field name
//TODO embed field
type structType struct {
	commonType
	fields    []structField //here need a stable order.
	fieldsMap map[string]*structField
}
type structField struct {
	Name string
	Type typeInterface
	*reflect.StructField
}

func (t *structType) init() {
	if t.fields != nil {
		return
	}
	t.fieldsMap = map[string]*structField{}
	for _, v := range kmgReflect.StructGetAllField(t.getReflectType()) {
		sf := structField{
			Name:        v.Name,
			Type:        t.ctx.mustNewTypeFromReflect(v.Type),
			StructField: v,
		}
		t.fields = append(t.fields, sf)
		t.fieldsMap[v.Name] = &sf
	}
}
func (t *structType) Html(v reflect.Value) template.HTML {
	t.init()
	type templateRow struct {
		Path string
		Name string
		Html template.HTML
	}
	var templateData []templateRow
	for _, sf := range t.fields {
		templateData = append(templateData, templateRow{
			Path: sf.Name,
			Name: sf.Name,
			Html: sf.Type.Html(v.FieldByIndex(sf.Index)),
		})
	}
	return theTemplate.MustExecuteNameToHtml("Struct", templateData)
}
func (t *structType) getSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	ev := v.FieldByName(k)
	if !ev.IsValid() {
		return reflect.Value{}, fmt.Errorf("field %s not find in struct", k)
	}
	return ev, nil
}

func (t *structType) Save(v *reflect.Value, path Path, value string) error {
	t.init()
	if len(path) == 0 {
		return fmt.Errorf("[structType.save] get struct with no path,value:%s", path, value)
	}

	ev := v.FieldByName(path[0])
	if !ev.IsValid() {
		return fmt.Errorf("[structType.save] field not find in struct,path:%s value:%s", path, value)
	}
	pEv := &ev
	err := t.fieldsMap[path[0]].Type.Save(pEv, path[1:], value)
	if err != nil {
		return err
	}

	//not change this struct
	if pEv == &ev {
		return nil
	}
	if v.CanSet() {
		return nil
	}
	output := reflect.New(t.getReflectType()).Elem()
	output.Set(*v)
	*v = output
	ev = v.FieldByName(path[0])
	ev.Set(*pEv)
	return nil
}
