package kmgReflect

import (
	"fmt"
	"reflect"
)

//it is func decl with reflect func
//if a field not have a name
//only actual params ,not include recv
type ReflectFuncDecl struct {
	Params         []*ReflectFuncDeclField
	Results        []*ReflectFuncDeclField
	ParamMap       map[string]*ReflectFuncDeclField
	ResultMap      map[string]*ReflectFuncDeclField
	ResultHasNames bool
}

type ReflectFuncDeclField struct {
	Name  string //may not have name
	Type  reflect.Type
	Index int //place in origin param or result ,method include recv ,start from 0
}

//reflectType may not match funcDecl
func NewReflectFuncDecl(funcDecl *FuncDecl, rFunc reflect.Type, isMethod bool) (*ReflectFuncDecl, error) {
	decl := &ReflectFuncDecl{
		ParamMap:  make(map[string]*ReflectFuncDeclField),
		ResultMap: make(map[string]*ReflectFuncDeclField),
	}
	paramIndex := 0
	if isMethod {
		paramIndex = 1
	}
	if len(funcDecl.ParamNames) != rFunc.NumIn()-paramIndex {
		return nil, fmt.Errorf("NewReflectFuncDecl(): Param num not match FuncDecl: %d,reflectType: %d",
			len(funcDecl.ParamNames),
			rFunc.NumIn()-paramIndex,
		)
	}
	for _, name := range funcDecl.ParamNames {
		field := &ReflectFuncDeclField{Name: name, Type: rFunc.In(paramIndex), Index: paramIndex}
		decl.Params = append(decl.Params, field)
		paramIndex++
	}
	//result not have name
	if len(funcDecl.ResultName) == 0 {
		decl.ResultHasNames = false
		for i := 0; i < rFunc.NumOut(); i++ {
			decl.Results = append(decl.Results, &ReflectFuncDeclField{Type: rFunc.In(paramIndex), Index: i})
		}
	} else {
		decl.ResultHasNames = true
		if len(funcDecl.ResultName) != rFunc.NumOut() {
			return nil, fmt.Errorf("NewReflectFuncDecl(): Result num not match FuncDecl: %d,reflectType: %d",
				len(funcDecl.ResultName),
				rFunc.NumOut(),
			)
		}
		resultIndex := 0
		//result have name
		for _, name := range funcDecl.ResultName {
			field := &ReflectFuncDeclField{Name: name, Type: rFunc.Out(resultIndex), Index: resultIndex}
			decl.Results = append(decl.Results, field)
			resultIndex++
		}
	}
	for _, field := range decl.Params {
		decl.ParamMap[field.Name] = field
	}
	if decl.ResultHasNames {
		for _, field := range decl.Results {
			decl.ResultMap[field.Name] = field
		}
	}
	return decl, nil
}
