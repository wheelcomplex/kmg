package kmgReflect

import (
	"go/ast"
	"go/doc"
	//"fmt"
)

type PackageDecl struct {
	TypeMap map[string]*TypeDecl
}

func NewPackageDeclFromAstPackage(astPkg *ast.Package, fullImportPath string) *PackageDecl {
	output := &PackageDecl{TypeMap: make(map[string]*TypeDecl)}

	docPkg := doc.New(astPkg, fullImportPath, doc.AllMethods)
	for _, t := range docPkg.Types {
		output.TypeMap[t.Name] = NewTypeDeclFromDocType(t)
	}
	return output
}
