package builder

import (
	"go/ast"
	"go/token"
	"go/types"

	. "github.com/dave/jennifer/jen"
)

type PkgFile struct {
	astFile  *ast.File
	fset     *token.FileSet
	gendecls []*ast.GenDecl
	FileName string
	PkgName  string
	pkgScope *types.Scope
}

func (file PkgFile) GenerateBuilder() string {
	f := NewFile(file.PkgName)

	structs := file.parsePkgStructs()
	for _, st := range structs {
		st.DefineBuilderStruct(f)
		st.DefineBuilderInitializer(f)
		st.DefineBuilderConstructors(f)
		st.DefineBuildFunc(f)
	}

	return f.GoString()
}

func (file PkgFile) parsePkgStructs() (pkgStructs []PkgStruct) {
	for _, decl := range file.gendecls {
		for _, spec := range decl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			st := file.pkgScope.Lookup(typeSpec.Name.Name)
			sturctMeta := st.Type().Underlying().(*types.Struct)

			pkgStruct := PkgStruct{
				fset: file.fset,
				name: typeSpec.Name.Name,
				meta: sturctMeta,
			}
			pkgStructs = append(pkgStructs, pkgStruct)
		}
	}

	return
}
