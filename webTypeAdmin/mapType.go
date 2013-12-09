package webTypeAdmin

import (
	"fmt"
	"html/template"
	"reflect"
	//"strconv"
)

//path -> key(Key type)
//key can be bool,string,stringEnum,int,float,
type mapType struct {
	commonType
	keyType            typeInterface
	keyStringConverter stringConverterType
	elemType           typeInterface
}

//key is not string?
func (t *mapType) init() {
	if t.elemType != nil {
		return
	}
	t.keyType = t.ctx.mustNewTypeFromReflect(t.getReflectType().Key())
	var ok bool
	t.keyStringConverter, ok = t.keyType.(stringConverterType)
	if !ok {
		panic(fmt.Sprintf(
			"mapType key type not implement stringConverterType,key: %s",
			t.keyType.getReflectType().Kind().String(),
		))
	}
	t.elemType = t.ctx.mustNewTypeFromReflect(t.getReflectType().Elem())
}
func (t *mapType) Html(v reflect.Value) template.HTML {
	t.init()
	type templateRow struct {
		Path string
		Key  string
		Html template.HTML
	}
	var templateData []templateRow
	for _, key := range v.MapKeys() {
		keyS := t.keyStringConverter.toString(key)
		val := v.MapIndex(key)
		templateData = append(templateData, templateRow{
			Path: keyS,
			Key:  keyS,
			Html: t.elemType.Html(val),
		})
	}
	return theTemplate.MustExecuteNameToHtml("Map", templateData)
}

func (t *mapType) getSubValueByString(v reflect.Value, k string) (reflect.Value, error) {
	t.init()
	vk, err := t.keyStringConverter.fromString(k)
	if err != nil {
		return reflect.Value{}, err
	}
	vv := v.MapIndex(vk)
	if !vv.IsValid() {
		return reflect.Value{}, fmt.Errorf("[mapType.getSubValueByString] map key not found k:%s", k)
	}
	return vv, nil
}
func (t *mapType) delete(v reflect.Value, k string) error {
	t.init()
	vk, err := t.keyStringConverter.fromString(k)
	if err != nil {
		return err
	}
	v.SetMapIndex(vk, reflect.Value{})
	return nil
}
func (t *mapType) create(v reflect.Value, k string) error {
	t.init()
	if v.IsNil() {
		v.Set(reflect.MakeMap(t.getReflectType()))
	}
	vk, err := t.keyStringConverter.fromString(k)
	if err != nil {
		return err
	}
	v.SetMapIndex(vk, reflect.New(t.elemType.getReflectType()).Elem())
	return nil
}

// MapIndex return a not addressable reflect.Value problem..
func (t *mapType) mapSave(m reflect.Value, k string, v string) error {
	t.init()
	if m.IsNil() {
		return fmt.Errorf("[mapType.mapSave] k: %s v: %s", k, v)
	}
	vk, err := t.keyStringConverter.fromString(k)
	if err != nil {
		return err
	}
	ev := reflect.New(t.elemType.getReflectType()).Elem()
	err = t.elemType.save(ev, v)
	if err != nil {
		return err
	}
	m.SetMapIndex(vk, ev)
	return nil
}

func (t *mapType) Save(v *reflect.Value, path Path, value string) error {
	if len(path) == 0 {
		return fmt.Errorf("[mapType.save] get map with no path, value:%s", value)
	}
	t.init()
	OriginCanSet := v.CanSet()
	if v.IsNil() {
		if OriginCanSet {
			v.Set(reflect.MakeMap(t.getReflectType()))
		} else {
			*v = reflect.MakeMap(t.getReflectType())
		}
	} else {
		//copy all exist data,if this value can not set
		if !OriginCanSet {
			output := reflect.MakeMap(t.getReflectType())
			output.Set(*v)
			*v = output
		}
	}
	vk, err := t.keyStringConverter.fromString(path[0])
	if err != nil {
		return err
	}
	saveElemV := v.MapIndex(vk)
	KeyNotExist := false
	if !saveElemV.IsValid() {
		saveElemV = reflect.New(t.elemType.getReflectType()).Elem()
		KeyNotExist = true
	}
	elemV := &saveElemV
	err = t.elemType.Save(elemV, path[1:], value)
	if err != nil {
		return err
	}
	if KeyNotExist {
		v.SetMapIndex(vk, saveElemV)
	}
	if elemV != &saveElemV {
		v.SetMapIndex(vk, *elemV)
	}
	return nil
}
