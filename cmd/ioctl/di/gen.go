/*
 * Copyright © 2022-present the iocgo authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package di

import (
	"go/ast"
	"go/format"
	"strings"
	"time"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type IGenerator interface {
	RegisterMarkers(into *markers.Registry) error
	Generate(*genall.GenerationContext) error
}

var (
	PackageMarker                       = markers.Must(markers.MakeDefinition("ioc:package", markers.DescribesPackage, struct{}{}))
	InitRegisterMarker                  = markers.Must(markers.MakeDefinition("ioc:init-register", markers.DescribesPackage, struct{}{}))
	StarterMarker                       = markers.Must(markers.MakeDefinition("ioc:starter", markers.DescribesPackage, struct{}{}))
	ComponentMarker                     = markers.Must(markers.MakeDefinition("ioc:component", markers.DescribesType, struct{}{}))
	ComponentNameMarker                 = markers.Must(markers.MakeDefinition("ioc:component:name", markers.DescribesType, ""))
	ConfigurationMarker                 = markers.Must(markers.MakeDefinition("ioc:configuration", markers.DescribesType, struct{}{}))
	BeanMarker                          = markers.Must(markers.MakeDefinition("ioc:bean", markers.DescribesType, struct{}{}))
	BeanNameMarker                      = markers.Must(markers.MakeDefinition("ioc:bean:name", markers.DescribesType, ""))
	ScopeMarker                         = markers.Must(markers.MakeDefinition("ioc:scope", markers.DescribesType, ""))
	PrimaryMarker                       = markers.Must(markers.MakeDefinition("ioc:primary", markers.DescribesType, struct{}{}))
	LazyMarker                          = markers.Must(markers.MakeDefinition("ioc:lazy", markers.DescribesType, struct{}{}))
	ExposeMarker                        = markers.Must(markers.MakeDefinition("ioc:expose", markers.DescribesType, ""))
	ConfigurationPropertiesMarker       = markers.Must(markers.MakeDefinition("ioc:configuration-properties", markers.DescribesType, struct{}{}))
	ConfigurationPropertiesPrefixMarker = markers.Must(markers.MakeDefinition("ioc:configuration-properties:prefix", markers.DescribesType, ""))
	BindPrefixMarker                    = markers.Must(markers.MakeDefinition("ioc:bind:prefix", markers.DescribesType, ""))
	ConstructorMarker                   = markers.Must(markers.MakeDefinition("ioc:constructor", markers.DescribesType, ""))
	InitMarker                          = markers.Must(markers.MakeDefinition("ioc:init", markers.DescribesType, ""))
	DestroyMarker                       = markers.Must(markers.MakeDefinition("ioc:destroy", markers.DescribesType, ""))
	InjectMarker                        = markers.Must(markers.MakeDefinition("ioc:inject", markers.DescribesField, struct{}{}))
	InjectNameMarker                    = markers.Must(markers.MakeDefinition("ioc:inject:name", markers.DescribesField, ""))
	InjectNamedMarker                   = markers.Must(markers.MakeDefinition("ioc:inject:named", markers.DescribesType, ""))
)

var generatorMarkers = []*markers.Definition{
	PackageMarker,
	InitRegisterMarker,
	StarterMarker,
	ComponentMarker,
	ComponentNameMarker,
	ConfigurationMarker,
	BeanMarker,
	BeanNameMarker,
	ScopeMarker,
	PrimaryMarker,
	LazyMarker,
	ExposeMarker,
	ConfigurationPropertiesMarker,
	ConfigurationPropertiesPrefixMarker,
	BindPrefixMarker,
	ConstructorMarker,
	InitMarker,
	DestroyMarker,
	InjectMarker,
	InjectNameMarker,
	InjectNamedMarker,
}

type Generator struct {
	HeaderFile string `marker:",optional"`
	Year       string `marker:",optional"`
}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool { return true }
}

func (Generator) RegisterMarkers(into *markers.Registry) error {
	for _, def := range generatorMarkers {
		if err := into.Register(def); err != nil {
			return err
		}
	}
	return nil
}

func (g Generator) Generate(ctx *genall.GenerationContext) error {
	var headerText string
	generatedAt := time.Now().UTC().Format(time.RFC3339)
	if g.HeaderFile != "" {
		headerBytes, err := ctx.ReadFile(g.HeaderFile)
		if err != nil {
			return err
		}
		headerText = string(headerBytes)
	}

	iocCtx := NewIocGoGenerationContext(ctx, headerText, generatedAt)
	for _, root := range ctx.Roots {
		outContents := iocCtx.handlePackage(root)
		if len(outContents) == 0 {
			continue
		}
		WriteTo(ctx, root, outContents, "zz_generated.iocgo.go")
	}
	return nil
}

type IocGoGenerationContext struct {
	Collector   *markers.Collector
	Checker     *loader.TypeChecker
	HeaderText  string
	GeneratedAt string
}

func NewIocGoGenerationContext(ctx *genall.GenerationContext, headerText, generatedAt string) IocGoGenerationContext {
	return IocGoGenerationContext{Collector: ctx.Collector, Checker: ctx.Checker, HeaderText: headerText, GeneratedAt: generatedAt}
}

func (ctx *IocGoGenerationContext) handlePackage(root *loader.Package) []byte {
	root.NeedSyntax()

	_, err := ctx.Collector.MarkersInPackage(root)
	if err != nil {
		root.AddError(err)
		return nil
	}
	packageMarkers, err := markers.PackageMarkers(ctx.Collector, root)
	if err != nil {
		root.AddError(err)
		return nil
	}

	typeInfos := make([]*markers.TypeInfo, 0)
	if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		typeInfos = append(typeInfos, info)
	}); err != nil {
		root.AddError(err)
		return nil
	}

	functions := collectFunctionInfos(root, ctx.Collector.Registry)
	if !shouldGenerate(packageMarkers, typeInfos, functions) {
		return nil
	}

	ctx.Checker.Check(root)
	root.NeedTypesInfo()

	builder := newRuntimeCodeBuilder(root, packageMarkers, typeInfos, functions, ctx.HeaderText, ctx.GeneratedAt)
	outBytes, err := builder.Build()
	if err != nil {
		root.AddError(err)
		return nil
	}
	formatted, err := format.Source(outBytes)
	if err != nil {
		root.AddError(err)
		return outBytes
	}
	return formatted
}

type functionInfo struct {
	File    *ast.File
	Decl    *ast.FuncDecl
	Markers markers.MarkerValues
}

func collectFunctionInfos(root *loader.Package, registry *markers.Registry) []functionInfo {
	functions := make([]functionInfo, 0)
	for _, file := range root.Syntax {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			functions = append(functions, functionInfo{File: file, Decl: fn, Markers: parseFunctionMarkers(registry, fn.Doc)})
		}
	}
	return functions
}

func parseFunctionMarkers(registry *markers.Registry, doc *ast.CommentGroup) markers.MarkerValues {
	values := make(markers.MarkerValues)
	if registry == nil || doc == nil {
		return values
	}
	for _, comment := range doc.List {
		if !strings.HasPrefix(comment.Text, "//") {
			continue
		}
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if !strings.HasPrefix(text, "+") {
			continue
		}
		def := registry.Lookup(text, markers.DescribesType)
		if def == nil {
			continue
		}
		value, err := def.Parse(text)
		if err != nil {
			continue
		}
		values[def.Name] = append(values[def.Name], value)
	}
	return values
}

func shouldGenerate(packageMarkers markers.MarkerValues, typeInfos []*markers.TypeInfo, functions []functionInfo) bool {
	if markerExists(packageMarkers, InitRegisterMarker.Name) || markerExists(packageMarkers, StarterMarker.Name) || markerExists(packageMarkers, PackageMarker.Name) {
		return true
	}
	for _, info := range typeInfos {
		if markerExists(info.Markers, ComponentMarker.Name) || markerExists(info.Markers, ConfigurationMarker.Name) {
			return true
		}
	}
	for _, fn := range functions {
		if markerExists(fn.Markers, BeanMarker.Name) {
			return true
		}
	}
	return false
}

func markerExists(values markers.MarkerValues, name string) bool {
	items, ok := values[name]
	return ok && len(items) > 0
}

func markerString(values markers.MarkerValues, name string) string {
	items, ok := values[name]
	if !ok || len(items) == 0 || items[0] == nil {
		return ""
	}
	if value, ok := items[0].(string); ok {
		return value
	}
	return ""
}
