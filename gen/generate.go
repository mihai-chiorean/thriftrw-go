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
	"github.com/thriftrw/thriftrw-go/internal/goast"
	"github.com/thriftrw/thriftrw-go/plugin/api"
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

	if o.NoRecurse {
		if err := generateModule(m, importer, o); err != nil {
			return err
		}
	} else {
		err := m.Walk(func(m *compile.Module) error {
			if err := generateModule(m, importer, o); err != nil {
				return generateError{Name: m.ThriftPath, Reason: err}
			}
			return nil
		})
		if err != nil {
			return err
		}
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

		for path, contents := range res.Files {
			if strings.Contains(path, "..") {
				// TODO(abg): not the place to do this
				return fmt.Errorf("a plugin is trying to access a parent directory: %v", path)
			}

			fullPath := filepath.Join(o.OutputDir, path)
			directory := filepath.Dir(fullPath)

			if path[len(path)-3:] == ".go" {
				var err error
				// TODO(abg): Need to take all files in this directory into account when
				// removing unused imports. See how goimports does it:
				// https://github.com/golang/tools/blob/0e9f43fcb67267967af8c15d7dc54b373e341d20/imports/fix.go#L77
				//
				// We should maybe run this transform on ALL Go files after everything
				// is generated.
				contents, err = goast.Reformat(nil, fullPath, contents, goast.RemoveUnusedImports)
				if err != nil {
					// TODO(abg): Need plugin name and possibly file contents here
					return fmt.Errorf("failed to reformat %q: %v", path, err)
				}
			}

			if err := os.MkdirAll(directory, 0755); err != nil {
				return fmt.Errorf("could not create directory %q: %v", directory, err)
			}

			if err := ioutil.WriteFile(fullPath, contents, 0644); err != nil {
				return fmt.Errorf("failed to write %q: %v", fullPath, err)
			}
		}
	}

	return nil
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

func buildModuleMap(i thriftPackageImporter, m *compile.Module) (map[string]api.ModuleID, map[api.ModuleID]*api.Module, error) {
	pathToID := make(map[string]api.ModuleID)
	idToModule := make(map[api.ModuleID]*api.Module)
	nextID := api.ModuleID(1)

	err := m.Walk(func(m *compile.Module) error {
		id := nextID
		nextID++

		importPath, err := i.Package(m.ThriftPath)
		if err != nil {
			return err
		}

		dir, err := i.RelativePackage(m.ThriftPath)
		if err != nil {
			return err
		}

		pathToID[m.ThriftPath] = id
		idToModule[id] = &api.Module{
			Package:   importPath,
			Directory: dir,
		}
		return nil
	})

	return pathToID, idToModule, err
}

func buildGenerateRequest(i thriftPackageImporter, m *compile.Module) (*api.GenerateRequest, error) {
	type key struct{ ThriftPath, ServiceName string }

	moduleIDs, modules, err := buildModuleMap(i, m)
	if err != nil {
		return nil, err
	}

	nextServiceID := api.ServiceID(1)
	serviceIDs := make(map[key]api.ServiceID)
	services := make(map[api.ServiceID]*api.Service)
	var rootServices []api.ServiceID

	var buildService func(spec *compile.ServiceSpec) (api.ServiceID, error)
	buildService = func(spec *compile.ServiceSpec) (api.ServiceID, error) {
		var parentID *api.ServiceID
		if spec.Parent != nil {
			parent, err := buildService(spec.Parent)
			if err != nil {
				return 0, err
			}
			parentID = &parent
		}

		k := key{ThriftPath: spec.ThriftFile(), ServiceName: spec.Name}
		if id, ok := serviceIDs[k]; ok {
			return id, nil
		}

		id := nextServiceID
		nextServiceID++

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
			ParentID:  parentID,
			Functions: functions,
			ModuleID:  moduleIDs[spec.ThriftFile()],
		}
		serviceIDs[k] = id
		// TODO(abg): This should only be added if it belonged to the root module
		// if o.NoRecurse is set.
		rootServices = append(rootServices, id)
		return id, nil
	}

	for _, spec := range m.Services {
		if _, err := buildService(spec); err != nil {
			return nil, fmt.Errorf("failed to compile service %q: %v", spec.Name, err)
		}
	}

	// TODO(abg): Use o.NoRecurse to decide what the root services should be.

	return &api.GenerateRequest{
		RootServices: rootServices,
		Services:     services,
		Modules:      modules,
	}, nil
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
			Name:       goCase(spec.Name),
			ThriftName: spec.Name,
			Arguments:  args,
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
			Name: goCase(f.Name),
			Type: t,
		})
	}
	return args, nil
}

func simpleType(t api.SimpleType) *api.SimpleType {
	return &t
}

func buildType(i thriftPackageImporter, spec compile.TypeSpec, required bool) (*api.Type, error) {
	var t *api.Type
	switch spec {
	case compile.BoolSpec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeBool)}
	case compile.I8Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt8)}
	case compile.I16Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt16)}
	case compile.I32Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt32)}
	case compile.I64Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt64)}
	case compile.DoubleSpec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeFloat64)}
	case compile.StringSpec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeString)}
	case compile.BinarySpec:
		// Don't need to wrap into a ptr
		return &api.Type{SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)}}, nil
	default:
		// Not a primitive type. Try checking if it's a container.
	}

	if t != nil {
		if !required {
			t = &api.Type{PointerType: t}
		}
		return t, nil
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
			return &api.Type{KeyValueSliceType: &api.TypePair{Left: k, Right: v}}, nil
		}

		return &api.Type{MapType: &api.TypePair{Left: k, Right: v}}, nil

	case *compile.ListSpec:
		v, err := buildType(i, s.ValueSpec, true)
		if err != nil {
			return nil, err
		}

		return &api.Type{SliceType: v}, nil

	case *compile.SetSpec:
		v, err := buildType(i, s.ValueSpec, true)
		if err != nil {
			return nil, err
		}

		if !isHashable(s.ValueSpec) {
			return &api.Type{SliceType: v}, nil
		}

		return &api.Type{MapType: &api.TypePair{
			Left:  v,
			Right: &api.Type{SimpleType: simpleType(api.SimpleTypeStructEmpty)},
		}}, nil

	case *compile.EnumSpec:
		importPath, err := i.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		t = &api.Type{
			ReferenceType: &api.TypeReference{
				Name:    goCase(s.Name),
				Package: importPath,
			},
		}
		if !required {
			t = &api.Type{PointerType: t}
		}
		return t, nil

	case *compile.StructSpec:
		importPath, err := i.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		return &api.Type{
			PointerType: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:    goCase(s.Name),
					Package: importPath,
				},
			},
		}, nil

	case *compile.TypedefSpec:
		importPath, err := i.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		t = &api.Type{
			ReferenceType: &api.TypeReference{
				Name:    goCase(s.Name),
				Package: importPath,
			},
		}

		if !required && !isReferenceType(spec) {
			t = &api.Type{PointerType: t}
		}

		return t, nil
	default:
		panic(fmt.Sprintf("Unknown type (%T) %v", spec, spec))
	}
}
