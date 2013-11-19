package kmgReflect

import (
	"go/doc"
)

type TypeDecl struct {
	Methods map[string]*FuncDecl
}

func NewTypeDeclFromDocType(docType *doc.Type) *TypeDecl {
	decl := &TypeDecl{
		//	FullName: FullName,
		Methods: make(map[string]*FuncDecl),
	}
	for _, m := range docType.Methods {
		decl.Methods[m.Name] = NewFuncDeclFromDocFunc(m)
	}
	return decl
}
