package kmgReflect

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	//"fmt"
)

//全局定义(包含所有的package)
type ContextDecl struct {
	PackageMap map[string]*PackageDecl
}

func (ctxt *ContextDecl) GetTypeDeclByReflectType(reflectType reflect.Type) (*TypeDecl, bool) {
	reflectType = IndirectType(reflectType)
	pkg, ok := ctxt.PackageMap[reflectType.PkgPath()]
	if !ok {
		return nil, false
	}
	t, ok := pkg.TypeMap[reflectType.Name()]
	if !ok {
		return nil, false
	}
	return t, true
}
func (ctxt *ContextDecl) GetMethodDeclByReflectType(reflectType reflect.Type, methodName string) (*FuncDecl, bool) {
	t, ok := ctxt.GetTypeDeclByReflectType(reflectType)
	if !ok {
		return nil, false
	}
	f, ok := t.Methods[methodName]
	if !ok {
		return nil, false
	}
	return f, true
}

func NewContextDeclFromSrcPath(root string) (*ContextDecl, error) {
	decl := &ContextDecl{PackageMap: make(map[string]*PackageDecl)}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if info.Name()[0] == '.' {
			return filepath.SkipDir
		}
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments|parser.AllErrors)
		if err != nil {
			return err
		}
		pkg := pkgs[info.Name()]
		if pkg == nil {
			return nil
		}
		fullImportPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		fullImportPath = strings.Replace(fullImportPath, "\\", "/", -1)
		decl.PackageMap[fullImportPath] = NewPackageDeclFromAstPackage(pkg, fullImportPath)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return decl, nil
}
