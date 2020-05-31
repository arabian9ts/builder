package builder

import (
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	. "github.com/dave/jennifer/jen"
)

const (
	GETTER_TAG_VALUE = "get"
	SETTER_TAG_VALUE = "set"
	BUILD_TAG_VALUE  = "build"
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

		biilder, found := st.BuildTagValue(i)
		if found && biilder != "" {
			idef = biilder
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

func (st PkgStruct) BuildTagValue(fieldNum int) (buildname string, found bool) {
	tag := st.meta.Tag(fieldNum)
	buildname, found = reflect.StructTag(tag).Lookup(BUILD_TAG_VALUE)
	return
}

func (st PkgStruct) GetterTagValue(fieldNum int) (gettername string, found bool) {
	tag := st.meta.Tag(fieldNum)
	gettername, found = reflect.StructTag(tag).Lookup(GETTER_TAG_VALUE)
	return
}

func (st PkgStruct) SetterTagValue(fieldNum int) (settername string, found bool) {
	tag := st.meta.Tag(fieldNum)
	settername, found = reflect.StructTag(tag).Lookup(SETTER_TAG_VALUE)
	return
}

func (st PkgStruct) DefineAccessors(file *File) {
	// getter
	{
		receiver := strings.ToLower(st.name)
		for i := 0; i < st.meta.NumFields(); i++ {
			field := st.meta.Field(i)
			argType := field.Type().String()

			typeIdx := strings.LastIndex(argType, ".")
			if 0 < typeIdx {
				argType = argType[typeIdx+1:]
			}

			getter, found := st.GetterTagValue(i)
			if !found {
				continue
			}
			if getter == "" {
				getter = fmt.Sprintf("Get%s", strings.Title(field.Name()))
			}

			file.Func().Params(Id(receiver).Op("*").Id(st.name)).
				Id(getter).
				Params().
				Params(Id(argType)).
				Block(
					Return(Id(receiver).Op(".").Id(field.Name())),
				).
				Line()
		}
	}

	// setter
	{
		receiver := strings.ToLower(st.name)
		for i := 0; i < st.meta.NumFields(); i++ {
			field := st.meta.Field(i)
			argType := field.Type().String()
			argument := strings.ToLower(field.Name())

			typeIdx := strings.LastIndex(argType, ".")
			if 0 < typeIdx {
				argType = argType[typeIdx+1:]
			}

			setter, found := st.SetterTagValue(i)
			if !found {
				continue
			}
			if setter == "" {
				setter = fmt.Sprintf("Set%s", strings.Title(field.Name()))
			}

			file.Func().Params(Id(receiver).Op("*").Id(st.name)).
				Id(setter).
				Params(Id(argument).Id(argType)).
				Params().
				Block(
					Id(receiver).Op(".").Id(field.Name()).Op("=").Id(argument),
				).
				Line()
		}
	}
}
