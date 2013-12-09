package webTypeAdmin

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgType"
	"html/template"
	"reflect"
)

var fkRefReflectType = reflect.TypeOf((*FkRef)(nil)).Elem()

type FkRef interface {
	GetReferenceType() reflect.Type
}

//xxId refer to another object
//path -> pass to RefType
type fkRefType struct {
	kmgType.KmgType
	underlyingType     adminType
	referenceType      adminType
	referenceContainer reflect.Value
	keyStringConverter kmgType.StringConverterInterface
	ctx                *context
}

func (t *fkRefType) init() (err error) {
	if t.referenceType != nil {
		return
	}
	var ok bool
	rrt := reflect.Zero(t.GetReflectType()).Interface().(FkRef).GetReferenceType()
	t.referenceType = t.ctx.typeOfFromReflect(rrt)
	rc, ok := getFkRefContainerValue(t.ctx.rootValue, rrt)
	if !ok {
		return fmt.Errorf("[fkRefType.init] not found referenceContainer")
	}
	t.referenceContainer = rc
	if rc.Type().Key() != t.GetReflectType() {
		return fmt.Errorf("[fkRefType.init] Container key type is not same with self type")
	}
}

func getFkRefContainerValue(v reflect.Value, rrt reflect.Type) (reflect.Value, bool) {
	switch v.Kind() {
	case reflect.Map:
		if v.Type().Elem() == rrt {
			return v, true
		}
		return reflect.Value{}, false
	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Value{}, false
		}
		return getFkRefContainerValue(v.Elem(), rrt)
	case reflect.Array:
		//TODO check mutli container
		if v.IsNil() {
			return reflect.Value{}, false
		}
		return getFkRefContainerValue(v.Elem(), rrt)
	case reflect.Struct:
		//TODO check mutli container
		len := v.NumField()
		for i := 0; i < len; i++ {
			rv, ok := getFkRefContainerValue(v.Field(i), rrt)
			if ok {
				return rv, true
			}
		}
		return reflect.Value{}, false
	//not enter slice
	default:
		return reflect.Value{}, false
	}
	return reflect.Value{}, false
}
func (t *fkRefType) HtmlView(v reflect.Value) (html template.HTML, err error) {
	if err = t.init(); err != nil {
		return
	}
	var templateData selectTemplateData
	templateData.Value = t.keyStringConverter.ToString(v)
	for _, vk := range t.referenceContainer.MapKeys() {
		sk := t.keyStringConverter.ToString(vk)
		templateData.List = append(templateData.List, sk)
	}
	return theTemplate.ExecuteNameToHtml("Select", templateData)
}

func (t *fkRefType) SaveByPath(v *reflect.Value, path kmgType.Path, value string) (err error) {
	if err = t.KmgType.SaveByPath(v, path, value); err != nil {
		return
	}
	t.init()
	vk, err := t.keyStringConverter.FromString(value)
	if err != nil {
		return err
	}
	vv := t.referenceContainer.MapIndex(vk)
	if !vv.IsValid() {
		return fmt.Errorf("[fkRefType.save] save value not in container map:%s", value)
	}
	return t.underlyingType.save(v, value)
}
