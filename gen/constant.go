// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package gen

import (
	"fmt"
	"strconv"

	"github.com/thriftrw/thriftrw-go/compile"
)

// Constant generates code for `const` expressions in Thrift files.
func Constant(g Generator, c *compile.Constant) error {
	err := g.DeclareFromTemplate(
		`<if canBeConstant .Type>const<else>var<end> <constantName .Name> <typeReference .Type> = <constantValue .Value .Type>`,
		c,
		TemplateFunc("constantValue", ConstantValue),
		TemplateFunc("canBeConstant", canBeConstant),
		TemplateFunc("constantName", constantName),
	)
	return wrapGenerateError(c.Name, err)
}

// ConstantValue generates an expression containing the given constant value of
// the given type.
//
// The constant must already have been linked to the given type.
func ConstantValue(g Generator, c compile.ConstantValue, t compile.TypeSpec) (string, error) {
	switch v := c.(type) {
	case compile.ConstantBool:
		return constantBool(g, v, t)
	case compile.ConstantDouble:
		return constantDouble(g, v, t)
	case compile.ConstantInt:
		return constantInt(g, v, t)
	case compile.ConstantList:
		return constantList(g, v, t)
	case compile.ConstantMap:
		return constantMap(g, v, t)
	case compile.ConstantSet:
		return constantSet(g, v, t)
	case compile.ConstantString:
		return strconv.Quote(string(v)), nil
	case *compile.ConstantStruct:
		return constantStruct(g, v, t)
	case compile.EnumItemReference:
		return enumItemReference(g, v, t)
	case compile.ConstReference:
		return g.LookupConstantName(v.Target)
	default:
		panic(fmt.Sprintf("Unknown constant value %v (%T)", c, c))
	}
}

func castConstant(g Generator, t compile.TypeSpec, s string) (string, error) {
	n, err := typeName(g, t)
	if err != nil {
		return "", err
	}
	s = fmt.Sprintf("%v(%v)", n, s)
	return s, nil
}

func constantBool(g Generator, v compile.ConstantBool, t compile.TypeSpec) (_ string, err error) {
	s := "false"
	if v {
		s = "true"
	}
	if t != compile.BoolSpec {
		s, err = castConstant(g, t, s)
	}
	return s, err
}

func constantDouble(g Generator, v compile.ConstantDouble, t compile.TypeSpec) (_ string, err error) {
	s := fmt.Sprint(float64(v))
	if t != compile.DoubleSpec {
		s, err = castConstant(g, t, s)
	}
	return s, err
}

func constantInt(g Generator, v compile.ConstantInt, t compile.TypeSpec) (_ string, err error) {
	s := fmt.Sprint(int(v))
	switch t {
	case compile.I8Spec, compile.I16Spec, compile.I32Spec, compile.I64Spec:
		// do nothing
	default:
		s, err = castConstant(g, t, s)
	}
	return s, err
}

func constantList(g Generator, v compile.ConstantList, t compile.TypeSpec) (string, error) {
	valueSpec := compile.RootTypeSpec(t).(*compile.ListSpec).ValueSpec
	return g.TextTemplate(
		`
		<$valueType := .ValueSpec>
		<typeReference .Spec>{
			<range .Value>
				<constantValue . $valueType>,
			<end>
		}`, struct {
			Spec      compile.TypeSpec
			ValueSpec compile.TypeSpec
			Value     compile.ConstantList
		}{Spec: t, ValueSpec: valueSpec, Value: v},
		TemplateFunc("constantValue", ConstantValue))
}

func constantMap(g Generator, v compile.ConstantMap, t compile.TypeSpec) (string, error) {
	mapSpec := compile.RootTypeSpec(t).(*compile.MapSpec)
	keySpec := mapSpec.KeySpec
	valueSpec := mapSpec.ValueSpec
	return g.TextTemplate(
		`
		<$keyType := .KeySpec>
		<$valueType := .ValueSpec>
		<typeReference .Spec>{
			<range .Value>
				<if isHashable $keyType>
					<constantValue .Key $keyType>:
						<constantValue .Value $valueType>,
				<else>
					{
						Key: <constantValue .Key $keyType>,
						Value: <constantValue .Value $valueType>,
					},
				<end>
			<end>
		}`, struct {
			Spec      compile.TypeSpec
			KeySpec   compile.TypeSpec
			ValueSpec compile.TypeSpec
			Value     compile.ConstantMap
		}{Spec: t, KeySpec: keySpec, ValueSpec: valueSpec, Value: v},
		TemplateFunc("constantValue", ConstantValue))
}

func constantSet(g Generator, v compile.ConstantSet, t compile.TypeSpec) (string, error) {
	valueSpec := compile.RootTypeSpec(t).(*compile.SetSpec).ValueSpec
	return g.TextTemplate(
		`
		<$valueType := .ValueSpec>
		<typeReference .Spec>{
			<range .Value>
				<if isHashable $valueType>
					<constantValue . $valueType>: struct{}{},
				<else>
					<constantValue . $valueType>,
				<end>
			<end>
		}`, struct {
			Spec      compile.TypeSpec
			ValueSpec compile.TypeSpec
			Value     compile.ConstantSet
		}{Spec: t, ValueSpec: valueSpec, Value: v},
		TemplateFunc("constantValue", ConstantValue))
}

func constantStruct(g Generator, v *compile.ConstantStruct, t compile.TypeSpec) (string, error) {
	fields := compile.RootTypeSpec(t).(*compile.StructSpec).Fields
	return g.TextTemplate(
		`
		<$fields := .Fields>
		&<typeName .Spec>{
			<range $name, $value := .Value.Fields>
				<$field := $fields.FindByName $name>
				<if and (not $field.Required) (isPrimitiveType $field.Type)>
					<goCase $field.Name>: <constantValuePtr $value $field.Type>,
				<else>
					<goCase $field.Name>: <constantValue $value $field.Type>,
				<end>
			<end>
		}`, struct {
			Spec   compile.TypeSpec
			Fields compile.FieldGroup
			Value  *compile.ConstantStruct
		}{Spec: t, Fields: fields, Value: v},
		TemplateFunc("constantValue", ConstantValue),
		TemplateFunc("constantValuePtr", ConstantValuePtr),
	)
}

func enumItemReference(g Generator, v compile.EnumItemReference, t compile.TypeSpec) (_ string, err error) {
	s, err := g.TextTemplate(`<enumItemName .Enum .Item.Name>`, v, TemplateFunc("enumItemName", enumItemName))
	if err != nil {
		return "", err
	}
	if t != v.Enum {
		s, err = castConstant(g, t, s)
	}
	return s, err
}

// ConstantValuePtr generates an expression which is a pointer to a value of
// type $t.
func ConstantValuePtr(g Generator, c compile.ConstantValue, t compile.TypeSpec) (string, error) {
	if !isPrimitiveType(t) {
		return ConstantValue(g, c, t)
	}

	err := g.EnsureDeclared(
		`func _<.ThriftName>_ptr(v <typeReference .>) *<typeReference .> {
			return &v
		}`, t)
	if err != nil {
		return "", err
	}

	s, err := ConstantValue(g, c, t)
	s = fmt.Sprintf("_%v_ptr(%v)", t.ThriftName(), s)
	return s, err
}
