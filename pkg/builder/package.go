package builder

import (
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
)

type Package struct {
	fset    *token.FileSet
	astPkg  *ast.Package
	PkgName string
}

type FileLoadFilterFunc func(info os.FileInfo) bool

func (pkg *Package) packageFiles() map[string]*ast.File {
	if pkg.astPkg == nil {
		return make(map[string]*ast.File)
	}

	return pkg.astPkg.Files
}

func (pkg *Package) ParsePkgFiles() (files []PkgFile) {
	for _, f := range pkg.packageFiles() {
		gendecls := make([]*ast.GenDecl, 0, len(f.Decls))
		for _, decl := range f.Decls {
			gendecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			gendecls = append(gendecls, gendecl)
		}

		conf := types.Config{
			Importer: importer.Default(),
			Error: func(err error) {
				fmt.Printf("!!! %#v\n", err)
			},
		}

		// limit access over files in same package
		pkgMeta, err := conf.Check(f.Name.String(), pkg.fset, []*ast.File{f}, nil)
		if err != nil {
			fmt.Printf("!!! %#v\n", err)
			continue
		}

		file := PkgFile{
			astFile:  f,
			fset:     pkg.fset,
			gendecls: gendecls,
			FileName: pkg.fset.File(f.Pos()).Name(),
			PkgName:  f.Name.String(),
			pkgScope: pkgMeta.Scope(),
		}
		files = append(files, file)
	}

	return
}

func LoadPackage(pkgDir string, filter FileLoadFilterFunc) (pkg *Package, err error) {
	fset := token.NewFileSet()
	pkgm, err := parser.ParseDir(
		fset,
		filepath.FromSlash(pkgDir),
		filter,
		parser.ParseComments,
	)
	if err != nil {
		return
	}
	if len(pkgm) <= 0 {
		pkg = &Package{fset: fset}
		return
	}

	for k, v := range pkgm {
		if pkg != nil {
			err = errors.New("must be single package dir")
			return
		}

		pkg = &Package{
			fset:    fset,
			astPkg:  v,
			PkgName: k,
		}
	}

	return
}
