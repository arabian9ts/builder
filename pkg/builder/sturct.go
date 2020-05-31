package builder

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"

	. "github.com/dave/jennifer/jen"
)

type PkgStruct struct {
	fset *token.FileSet
	name string
	meta *types.Struct
}

func (st PkgStruct) receiverName() string {
	return fmt.Sprintf("%sBuilder", strings.ToLower(st.name))
}

func (st PkgStruct) builderName() string {
	return fmt.Sprintf("%sBuilder", strings.Title(st.name))
}

func (st PkgStruct) builderInitializerName() string {
	return fmt.Sprintf("New%sBuilder", strings.Title(st.name))
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
	fields := make([]Code, 0, st.meta.NumFields())
	for i := 0; i < st.meta.NumFields(); i++ {
		fld := st.meta.Field(i)
		if len(fld.Name()) <= 0 {
			continue
		}

		attrName := strings.ToLower(fld.Name())
		fieldName := fld.Type().String()

		typeIdx := strings.LastIndex(fieldName, ".")
		if 0 < typeIdx {
			fieldName = fieldName[typeIdx+1:]
		}

		field := Id(attrName).Id(fieldName)
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
	for i := 0; i < st.meta.NumFields(); i++ {
		field := st.meta.Field(i)
		if len(field.Name()) <= 0 {
			continue
		}

		argType := field.Type().String()
		attr := strings.ToLower(field.Name())
		idef := strings.Title(attr)
		argment := strings.ToLower(attr)

		typeIdx := strings.LastIndex(argType, ".")
		if 0 < typeIdx {
			argType = argType[typeIdx+1:]
		}

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
	for i := 0; i < st.meta.NumFields(); i++ {
		field := st.meta.Field(i)
		if len(field.Name()) <= 0 {
			continue
		}

		structAttr := field.Name()
		builderAttr := strings.ToLower(structAttr)
		dict[Id(structAttr)] = Id(receiver).Op(".").Id(builderAttr)
	}

	if len(dict) <= 0 {
		return
	}

	file.Func().Params(Id(receiver).Id(builder)).
		Id("Build").
		Params().
		Params(Op("*").Id(st.name)).
		Block(
			Return(
				Op("&").Id(st.name).Values(dict),
			),
		)
}
