package builder

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	. "github.com/dave/jennifer/jen"
)

type PkgStruct struct {
	fset          *token.FileSet
	astStructType *ast.StructType
	astFile       *ast.File
	StructName    string
	FileName      string
	PkgName       string
}

func (st PkgStruct) receiverName() string {
	return fmt.Sprintf("%sBuilder", strings.ToLower(st.StructName))
}

func (st PkgStruct) builderName() string {
	return fmt.Sprintf("%sBuilder", strings.Title(st.StructName))
}

func (st PkgStruct) builderInitializerName() string {
	return fmt.Sprintf("New%sBuilder", strings.Title(st.StructName))
}

func (st PkgStruct) DefineBuilderInitializer(file *File) {
	builder := st.builderName()
	initializer := st.builderInitializerName()
	file.Func().
		Id(initializer).Params().
		Params(Op("*").Id(builder)).
		Block(
			Return(
				Op("&").Id(builder).Block(),
			),
		).
		Line()
}

func (st PkgStruct) DefineBuilderStruct(file *File) {
	fields := make([]Code, 0, len(st.astStructType.Fields.List))
	for _, fld := range st.astStructType.Fields.List {
		if len(fld.Names) <= 0 {
			continue
		}

		var tbuf bytes.Buffer
		err := printer.Fprint(&tbuf, st.fset, fld.Type)
		if err != nil {
			continue
		}

		attr := strings.ToLower(fld.Names[0].Name)
		field := Id(attr).Id(tbuf.String())
		fields = append(fields, field)
	}

	if len(fields) <= 0 {
		return
	}

	builder := st.builderName()
	file.Type().Id(builder).Struct(fields...)
}

func (st PkgStruct) DefineBuilderConstructors(file *File) {
	builder := st.builderName()
	receiver := st.receiverName()
	for _, field := range st.astStructType.Fields.List {
		if len(field.Names) <= 0 {
			continue
		}

		var tbuf bytes.Buffer
		err := printer.Fprint(&tbuf, st.fset, field.Type)
		if err != nil {
			continue
		}

		argType := tbuf.String()
		attr := strings.ToLower(field.Names[0].Name)
		idef := strings.Title(attr)
		argment := strings.ToLower(attr)

		file.Func().Params(Id(receiver).Op("*").Id(builder)).
			Id(idef).
			Params(Id(argment).Id(argType)).
			Params(Op("*").Id(builder)).
			Block(
				Id(receiver).Op(".").Id(attr).Op("=").Id(strings.ToLower(attr)),
				Return(Id(receiver)),
			).
			Line()
	}
}

func (st PkgStruct) DefineBuildFunc(file *File) {
	dict := Dict{}
	builder := st.builderName()
	receiver := st.receiverName()
	for _, field := range st.astStructType.Fields.List {
		if len(field.Names) <= 0 {
			continue
		}

		var tbuf bytes.Buffer
		err := printer.Fprint(&tbuf, st.fset, field.Type)
		if err != nil {
			continue
		}

		structAttr := field.Names[0].Name
		builderAttr := strings.ToLower(structAttr)
		dict[Id(structAttr)] = Id(receiver).Op(".").Id(builderAttr)
	}

	if len(dict) <= 0 {
		return
	}

	file.Func().Params(Id(receiver).Id(builder)).
		Id("Build").
		Params().
		Params(Op("*").Id(st.StructName)).
		Block(
			Return(
				Op("&").Id(st.StructName).Values(dict),
			),
		)
}
