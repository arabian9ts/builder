package builder

import (
	"go/ast"
	"go/token"

	. "github.com/dave/jennifer/jen"
)

type PkgFile struct {
	astFile  *ast.File
	fset     *token.FileSet
	gendecls []*ast.GenDecl
	FileName string
	PkgName  string
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
			typeSpect, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			st, ok := typeSpect.Type.(*ast.StructType)
			if !ok {
				continue
			}

			pkgStruct := PkgStruct{
				fset:          file.fset,
				astStructType: st,
				astFile:       file.astFile,
				StructName:    typeSpect.Name.Name,
				FileName:      file.FileName,
				PkgName:       file.PkgName,
			}
			pkgStructs = append(pkgStructs, pkgStruct)
		}
	}

	return
}
