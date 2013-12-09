package webTypeAdmin

import (
	"fmt"
	"html/template"
	"reflect"
)

var fkRefReflectType = reflect.TypeOf((*FkRef)(nil)).Elem()

/*
define a FkRef like this
type UserId int
func (user_id UserId)GetReferenceType()reflect.Type{
	return reflect.TypeOf((*User)(nil)).Elem()
}
*/
type FkRef interface {
	GetReferenceType() reflect.Type
}

//xxId refer to another object
//path -> pass to RefType
type fkRefType struct {
	commonType
	underlyingType     typeInterface
	referenceType      typeInterface
	referenceContainer reflect.Value
	stringConverterType
}

func (t *fkRefType) init() {
	if t.referenceType != nil {
		return
	}
	var ok bool
	rrt := reflect.Zero(t.getReflectType()).Interface().(FkRef).GetReferenceType()
	t.referenceType = t.ctx.mustNewTypeFromReflect(rrt)
	rc, ok := getFkRefContainerValue(t.ctx.rootValue, rrt)
	if !ok {
		panic("[fkRefType.init] not found referenceContainer")
	}
	t.referenceContainer = rc
	if rc.Type().Key() != t.getReflectType() {
		panic("[fkRefType.init] Container key type is not same with self type")
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
func (t *fkRefType) Html(v reflect.Value) template.HTML {
	t.init()
	var templateData selectTemplateData
	templateData.Value = t.stringConverterType.toString(v)
	for _, vk := range t.referenceContainer.MapKeys() {
		sk := t.stringConverterType.toString(vk)
		templateData.List = append(templateData.List, sk)
	}
	return theTemplate.MustExecuteNameToHtml("Select", templateData)
}
func (t *fkRefType) save(v reflect.Value, value string) error {
	t.init()
	vk, err := t.stringConverterType.fromString(value)
	if err != nil {
		return err
	}
	vv := t.referenceContainer.MapIndex(vk)
	if !vv.IsValid() {
		return fmt.Errorf("[fkRefType.save] save value not in container map:%s", value)
	}
	return t.underlyingType.save(v, value)
}

func (t *fkRefType) Save(v *reflect.Value, path Path, value string) error {
	if err := scaleValueSaveHandle(t, &v, path); err != nil {
		return err
	}
	t.init()
	vk, err := t.stringConverterType.fromString(value)
	if err != nil {
		return err
	}
	vv := t.referenceContainer.MapIndex(vk)
	if !vv.IsValid() {
		return fmt.Errorf("[fkRefType.save] save value not in container map:%s", value)
	}
	return t.underlyingType.save(v, value)
}
