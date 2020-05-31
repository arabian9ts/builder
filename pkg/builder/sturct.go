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

type Field struct {
	tag string
	*types.Var
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
	if len(st.filterOpenedFields()) <= 0 {
		return
	}

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
	for _, fld := range st.filterOpenedFields() {
		if len(fld.Name()) <= 0 {
			continue
		}

		fieldName := fld.Type().String()
		typeIdx := strings.LastIndex(fieldName, ".")
		if 0 < typeIdx {
			fieldName = fieldName[typeIdx+1:]
		}

		field := Id(fld.Name()).Id(fieldName)
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
	for _, field := range st.filterOpenedFields() {
		if len(field.Name()) <= 0 {
			continue
		}

		argType := field.Type().String()
		idef := strings.Title(field.Name())
		argment := strings.ToLower(field.Name())

		typeIdx := strings.LastIndex(argType, ".")
		if 0 < typeIdx {
			argType = argType[typeIdx+1:]
		}

		build, _ := field.BuildTagValue()
		if build == "-" {
			continue
		}
		if build != "" {
			idef = build
		}

		file.Func().Params(Id(receiver).Op("*").Id(builder)).
			Id(idef).
			Params(Id(argment).Id(argType)).
			Params(Op("*").Id(builder)).
			Block(
				Id(receiver).Op(".").Id(field.Name()).Op("=").Id(strings.ToLower(field.Name())),
				Return(Id(receiver)),
			).
			Line()
	}
}

func (st PkgStruct) DefineBuildFunc(file *File) {
	dict := Dict{}
	builder := st.builderName()
	receiver := st.receiverName()
	for _, field := range st.filterOpenedFields() {
		if len(field.Name()) <= 0 {
			continue
		}

		dict[Id(field.Name())] = Id(receiver).Op(".").Id(field.Name())
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

func (f Field) BuildTagValue() (buildname string, found bool) {
	buildname, found = reflect.StructTag(f.tag).Lookup(BUILD_TAG_VALUE)
	return
}

func (f Field) GetterTagValue() (gettername string, found bool) {
	gettername, found = reflect.StructTag(f.tag).Lookup(GETTER_TAG_VALUE)
	return
}

func (f Field) SetterTagValue() (settername string, found bool) {
	settername, found = reflect.StructTag(f.tag).Lookup(SETTER_TAG_VALUE)
	return
}

func (st PkgStruct) DefineAccessors(file *File) {
	// getter
	{
		receiver := strings.ToLower(st.name)
		for _, field := range st.filterOpenedFields() {
			argType := field.Type().String()

			typeIdx := strings.LastIndex(argType, ".")
			if 0 < typeIdx {
				argType = argType[typeIdx+1:]
			}

			getter, found := field.GetterTagValue()
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
		for _, field := range st.filterOpenedFields() {
			argType := field.Type().String()
			argument := strings.ToLower(field.Name())

			typeIdx := strings.LastIndex(argType, ".")
			if 0 < typeIdx {
				argType = argType[typeIdx+1:]
			}

			setter, found := field.SetterTagValue()
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

func (st PkgStruct) filterOpenedFields() (fields []Field) {
	for i := 0; i < st.meta.NumFields(); i++ {
		field := st.meta.Field(i)
		if field.Name() == strings.Title(field.Name()) {
			continue
		}

		fields = append(fields, Field{st.meta.Tag(i), field})
	}

	return
}
