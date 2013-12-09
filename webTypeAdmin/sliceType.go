package webTypeAdmin

import (
	"fmt"
	"html/template"
	"reflect"
	//"strconv"
)

//path -> slice index(number)
type sliceType struct {
	commonType
	elemType typeInterface
}

func (t *sliceType) init() {
	if t.elemType != nil {
		return
	}
	t.elemType = t.ctx.mustNewTypeFromReflect(t.getReflectType().Elem())
}
func (t *sliceType) Html(v reflect.Value) template.HTML {
	t.init()
	type templateRow struct {
		Path  int
		Index int
		Html  template.HTML
	}
	var templateData []templateRow
	len := v.Len()
	for i := 0; i < len; i++ {
		templateData = append(templateData, templateRow{
			Path:  i,
			Index: i,
			Html:  t.elemType.Html(v.Index(i)),
		})
	}
	return theTemplate.MustExecuteNameToHtml("Slice", templateData)
}

func (t *sliceType) getSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	t.init()
	return arrayGetSubValueByString(v, k)
}
func (t *sliceType) delete(v reflect.Value, k string) error {
	t.init()
	i, err := arrayParseKey(v, k)
	if err != nil {
		return err
	}
	v.Set(
		reflect.AppendSlice(v.Slice(0, i), v.Slice(i+1, v.Len())),
	)
	return nil
}
func (t *sliceType) create(v reflect.Value, k string) error {
	t.init()
	v.Set(
		reflect.Append(v, reflect.New(t.elemType.getReflectType()).Elem()),
	)
	return nil
}

func (t *sliceType) Save(v *reflect.Value, path Path, value string) error {
	t.init()
	if len(path) == 0 {
		return fmt.Errorf("[sliceType.save] get struct with no path,value:%s", path, value)
	}

	i, err := arrayParseKey(v, path[0])
	if err != nil {
		return err
	}
	ev := v.Index(i)
	pEv := &ev
	err := t.elemType.Save(pEv, path[1:], value)
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
	ev = v.Index(i)
	ev.Set(*pEv)
	return nil
}
