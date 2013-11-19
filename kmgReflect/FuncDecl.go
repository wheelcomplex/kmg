package kmgReflect

import (
	"go/ast"
	"go/doc"
)

//golang missing function decl of parameter name
//it is just too complex to get type info...
type FuncDecl struct {
	ParamNames []string //paramNames
	ResultName []string //resultNames
}

func NewFuncDeclFromDocFunc(docfunc *doc.Func) *FuncDecl {
	decl := &FuncDecl{}
	funcType := docfunc.Decl.Type
	for _, p := range funcType.Params.List {
		decl.ParamNames = append(decl.ParamNames, getNameFromAstField(p)...)
	}
	if funcType.Results != nil {
		for _, p := range funcType.Results.List {
			decl.ResultName = append(decl.ResultName, getNameFromAstField(p)...)
		}
	}
	return decl
}

func getNameFromAstField(astField *ast.Field) []string {
	output := []string{}
	//some thing like "func test1(a,b int)"->astField->"a,b"
	for _, name := range astField.Names {
		output = append(output, name.Name)
	}
	return output
}
