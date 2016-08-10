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
	"bytes"
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/thriftrw/thriftrw-go/compile"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/ptr"
)

// Options controls how code gets generated.
type Options struct {
	// OutputDir is the directory into which all generated code is written.
	//
	// This must be an absolute path.
	OutputDir string

	// PackagePrefix controls the import path prefix for all generated
	// packages.
	PackagePrefix string

	// ThriftRoot is the directory within whose tree all Thrift files consumed
	// are contained. The locations of the Thrift files relative to the
	// ThriftFile determines the module structure in OutputDir.
	//
	// This must be an absolute path.
	ThriftRoot string

	// NoRecurse determines whether code should be generated for included Thrift
	// files as well. If true, code gets generated only for the first module.
	NoRecurse bool

	// List of plugins
	Plugins []Plug
}

// Generate generates code based on the given options.
func Generate(m *compile.Module, o *Options) error {
	if !filepath.IsAbs(o.ThriftRoot) {
		return fmt.Errorf(
			"ThriftRoot must be an absolute path: %q is not absolute",
			o.ThriftRoot)
	}

	if !filepath.IsAbs(o.OutputDir) {
		return fmt.Errorf(
			"OutputDir must be an absolute path: %q is not absolute",
			o.OutputDir)
	}

	importer := thriftPackageImporter{
		ImportPrefix: o.PackagePrefix,
		ThriftRoot:   o.ThriftRoot,
	}

	if len(o.Plugins) > 0 {
		plug := multiPlug(o.Plugins)
		if err := plug.Open(); err != nil {
			return err
		}
		defer plug.Close()

		req, err := buildGenerateRequest(importer, m)
		if err != nil {
			return err
		}

		res, err := plug.Generate(req)
		if err != nil {
			return err
		}

		// TODO(abg): generate files
		fmt.Println("Generating", res)
	}

	if o.NoRecurse {
		return generateModule(m, importer, o)
	}

	return m.Walk(func(m *compile.Module) error {
		if err := generateModule(m, importer, o); err != nil {
			return generateError{Name: m.ThriftPath, Reason: err}
		}
		return nil
	})
}

// TODO(abg): Make some sort of public interface out of the Importer

type thriftPackageImporter struct {
	ImportPrefix string
	ThriftRoot   string
}

// RelativePackage returns the import path for the top-level package of the
// given Thrift file relative to the ImportPrefix.
func (i thriftPackageImporter) RelativePackage(file string) (string, error) {
	return filepath.Rel(i.ThriftRoot, strings.TrimSuffix(file, ".thrift"))
}

// Package returns the import path for the top-level package of the given Thrift
// file.
func (i thriftPackageImporter) Package(file string) (string, error) {
	pkg, err := i.RelativePackage(file)
	if err != nil {
		return "", err
	}
	return filepath.Join(i.ImportPrefix, pkg), nil
}

// ServicePackage returns the import path for the package for the Thrift service
// with the given name defined in the given file.
func (i thriftPackageImporter) ServicePackage(file, name string) (string, error) {
	topPackage, err := i.Package(file)
	if err != nil {
		return "", err
	}

	return filepath.Join(topPackage, "service", strings.ToLower(name)), nil
}

// generates code for only the given module, assuming that code for included
// modules has already been generated.
func generateModule(m *compile.Module, i thriftPackageImporter, o *Options) error {
	outDir := o.OutputDir
	// packageRelPath is the path relative to outputDir into which we'll be
	// writing the package for this Thrift file. For $thriftRoot/foo/bar.thrift,
	// packageRelPath is foo/bar, and packageDir is $outputDir/foo/bar. All
	// files for bar.thrift will be written to the $outputDir/foo/bar/ tree. The
	// package will be importable via $importPrefix/foo/bar.
	packageRelPath, err := i.RelativePackage(m.ThriftPath)
	if err != nil {
		return err
	}

	// TODO(abg): Prefer top-level package name from `namespace go` directive.
	packageName := filepath.Base(packageRelPath)

	// importPath is the full import path for the top-level package generated
	// for this Thrift file.
	importPath, err := i.Package(m.ThriftPath)
	if err != nil {
		return err
	}

	// packageOutDir is the directory whithin which all files and folders for
	// this Thrift file will be written.
	packageOutDir := filepath.Join(outDir, packageRelPath)

	// Mapping of file names relative to $packageOutDir and their contents.
	files := make(map[string]*bytes.Buffer)

	if len(m.Constants) > 0 {
		g := NewGenerator(i, importPath, packageName)

		for _, constantName := range sortStringKeys(m.Constants) {
			if err := Constant(g, m.Constants[constantName]); err != nil {
				return err
			}
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, token.NewFileSet()); err != nil {
			return fmt.Errorf(
				"could not generate constants for %q: %v", m.ThriftPath, err)
		}

		// TODO(abg): Verify no file collisions
		files["constants.go"] = buff
	}

	if len(m.Types) > 0 {
		g := NewGenerator(i, importPath, packageName)

		for _, typeName := range sortStringKeys(m.Types) {
			if err := TypeDefinition(g, m.Types[typeName]); err != nil {
				return err
			}
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, token.NewFileSet()); err != nil {
			return fmt.Errorf(
				"could not generate types for %q: %v", m.ThriftPath, err)
		}

		// TODO(abg): Verify no file collisions
		files["types.go"] = buff
	}

	if len(m.Services) > 0 {
		for _, serviceName := range sortStringKeys(m.Services) {
			service := m.Services[serviceName]
			importPath, err := i.ServicePackage(service.ThriftFile(), service.Name)
			if err != nil {
				return err
			}
			packageName := filepath.Base(importPath)

			// TODO inherited service functions

			g := NewGenerator(i, importPath, packageName)
			serviceFiles, err := Service(g, service)
			if err != nil {
				return fmt.Errorf(
					"could not generate code for service %q: %v",
					serviceName, err)
			}

			for name, buff := range serviceFiles {
				filename := filepath.Join("service", packageName, name)
				files[filename] = buff
			}
		}
	}

	for relPath, contents := range files {
		fullPath := filepath.Join(packageOutDir, relPath)
		directory := filepath.Dir(fullPath)

		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("could not create directory %q: %v", directory, err)
		}

		if err := ioutil.WriteFile(fullPath, contents.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write %q: %v", fullPath, err)
		}
	}
	return nil
}

func buildGenerateRequest(i thriftPackageImporter, m *compile.Module) (*api.GenerateRequest, error) {
	type key struct{ ThriftPath, ServiceName string }

	nextID := int32(1)
	serviceIds := make(map[key]int32)
	services := make(map[int32]*api.Service)

	var buildService func(spec *compile.ServiceSpec) (int32, error)
	buildService = func(spec *compile.ServiceSpec) (int32, error) {
		var parentID *int32
		if spec.Parent != nil {
			parent, err := buildService(spec.Parent)
			if err != nil {
				return 0, err
			}
			parentID = &parent
		}

		k := key{ThriftPath: spec.ThriftFile(), ServiceName: spec.Name}
		if id, ok := serviceIds[k]; ok {
			return id, nil
		}

		id := nextID
		nextID++

		importPath, err := i.ServicePackage(spec.ThriftFile(), spec.Name)
		if err != nil {
			return 0, err
		}

		dir, err := i.RelativePackage(spec.ThriftFile())
		if err != nil {
			return 0, err
		}

		functions, err := buildFunctions(i, spec.Functions)
		if err != nil {
			return 0, err
		}

		// TODO make this part of generateModule somehow
		services[id] = &api.Service{
			Name:      spec.Name,
			Package:   importPath,
			Directory: filepath.Join(dir, "service", filepath.Base(importPath)),
			ParentId:  parentID,
			Functions: functions,
		}
		serviceIds[k] = id
		return id, nil
	}

	for _, spec := range m.Services {
		if _, err := buildService(spec); err != nil {
			return nil, fmt.Errorf("failed to compile service %q: %v", spec.Name, err)
		}
	}

	return &api.GenerateRequest{Services: services}, nil
}

func buildFunctions(i thriftPackageImporter, fs map[string]*compile.FunctionSpec) ([]*api.Function, error) {
	functions := make([]*api.Function, 0, len(fs))
	for _, functionName := range sortStringKeys(fs) {
		spec := fs[functionName]

		args, err := buildArguments(i, compile.FieldGroup(spec.ArgsSpec))
		if err != nil {
			return nil, err
		}
		function := &api.Function{
			Name:      spec.Name,
			Arguments: args,
		}

		if spec.ResultSpec != nil {
			var err error
			result := spec.ResultSpec
			if result.ReturnType != nil {
				function.ReturnType, err = buildType(i, result.ReturnType, true)
				if err != nil {
					return nil, err
				}
			}
			if len(result.Exceptions) > 0 {
				function.Exceptions, err = buildArguments(i, result.Exceptions)
				if err != nil {
					return nil, err
				}
			}
		}

		functions = append(functions, function)
	}
	return functions, nil
}

func buildArguments(i thriftPackageImporter, fs compile.FieldGroup) ([]*api.Argument, error) {
	args := make([]*api.Argument, 0, len(fs))
	for _, f := range fs {
		t, err := buildType(i, f.Type, f.Required)
		if err != nil {
			return nil, err
		}

		args = append(args, &api.Argument{
			Name: f.Name, // TODO goCase
			Type: t,
		})
	}
	return args, nil
}

func buildType(i thriftPackageImporter, spec compile.TypeSpec, required bool) (*api.Type, error) {
	simpleType := func(t api.SimpleType) *api.SimpleType {
		return &t
	}

	switch spec {
	case compile.BoolSpec:
		return &api.Type{
			Info:    &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeBool)},
			Pointer: ptr.Bool(!required),
		}, nil
	case compile.I8Spec:
		return &api.Type{
			Info:    &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeInt8)},
			Pointer: ptr.Bool(!required),
		}, nil
	case compile.I16Spec:
		return &api.Type{
			Info:    &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeInt16)},
			Pointer: ptr.Bool(!required),
		}, nil
	case compile.I32Spec:
		return &api.Type{
			Info:    &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeInt32)},
			Pointer: ptr.Bool(!required),
		}, nil
	case compile.I64Spec:
		return &api.Type{
			Info:    &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeInt64)},
			Pointer: ptr.Bool(!required),
		}, nil
	case compile.DoubleSpec:
		return &api.Type{
			Info:    &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeFloat64)},
			Pointer: ptr.Bool(!required),
		}, nil
	case compile.StringSpec:
		return &api.Type{
			Info:    &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeString)},
			Pointer: ptr.Bool(!required),
		}, nil
	case compile.BinarySpec:
		return &api.Type{
			Info: &api.TypeInfo{
				SliceType: &api.SliceType{
					ValueType: &api.Type{Info: &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeByte)}},
				},
			},
		}, nil
	default:
		// Not a primitive type. Try checking if it's a container.
	}

	switch s := spec.(type) {
	case *compile.MapSpec:
		k, err := buildType(i, s.KeySpec, true)
		if err != nil {
			return nil, err
		}

		v, err := buildType(i, s.ValueSpec, true)
		if err != nil {
			return nil, err
		}

		if !isHashable(s.KeySpec) {
			return &api.Type{
				Info: &api.TypeInfo{
					KeyValueSliceType: &api.KeyValueSliceType{KeyType: k, ValueType: v},
				},
			}, nil
		}

		return &api.Type{
			Info: &api.TypeInfo{
				MapType: &api.MapType{KeyType: k, ValueType: v},
			},
		}, nil

	case *compile.ListSpec:
		v, err := buildType(i, s.ValueSpec, true)
		if err != nil {
			return nil, err
		}
		return &api.Type{
			Info: &api.TypeInfo{
				SliceType: &api.SliceType{ValueType: v},
			},
		}, nil

	case *compile.SetSpec:
		v, err := buildType(i, s.ValueSpec, true)
		if err != nil {
			return nil, err
		}

		if !isHashable(s.ValueSpec) {
			return &api.Type{
				Info: &api.TypeInfo{
					SliceType: &api.SliceType{ValueType: v},
				},
			}, nil
		}
		return &api.Type{
			Info: &api.TypeInfo{
				MapType: &api.MapType{
					KeyType: v,
					ValueType: &api.Type{
						Info: &api.TypeInfo{SimpleType: simpleType(api.SimpleTypeStructEmpty)},
					},
				},
			},
		}, nil

	case *compile.EnumSpec:
		importPath, err := i.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		return &api.Type{
			Info: &api.TypeInfo{
				ReferenceType: &api.TypeReference{
					Name:    s.Name,
					Package: importPath,
				},
			},
			Pointer: ptr.Bool(!required),
		}, nil

	case *compile.StructSpec:
		importPath, err := i.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		return &api.Type{
			Info: &api.TypeInfo{
				ReferenceType: &api.TypeReference{
					Name:    s.Name,
					Package: importPath,
				},
			},
			Pointer: ptr.Bool(true),
		}, nil

	case *compile.TypedefSpec:
		importPath, err := i.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		return &api.Type{
			Info: &api.TypeInfo{
				ReferenceType: &api.TypeReference{
					Name:    s.Name,
					Package: importPath,
				},
			},
			Pointer: ptr.Bool(!required && !isReferenceType(spec)),
		}, nil
	default:
		panic(fmt.Sprintf("Unknown type (%T) %v", spec, spec))
	}
}
