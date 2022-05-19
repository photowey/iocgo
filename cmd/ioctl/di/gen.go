/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package di

import (
	"bytes"
	"go/ast"
	"go/format"
	"io"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type IGenerator interface {
	RegisterMarkers(into *markers.Registry) error
	Generate(*genall.GenerationContext) error
}

var (
	AutowiredMarker     = markers.Must(markers.MakeDefinition("ioc:autowired", markers.DescribesType, true))             // autowired: default: true
	ScopeMarker         = markers.Must(markers.MakeDefinition("ioc:autowired:scope", markers.DescribesType, ""))         // scope: singleton, prototype
	InterfacesMarker    = markers.Must(markers.MakeDefinition("ioc:autowired:interfaces", markers.DescribesType, ""))    // interfaces
	ParameterMarker     = markers.Must(markers.MakeDefinition("ioc:autowired:parameter", markers.DescribesType, ""))     // parameter
	ConstructorMarker   = markers.Must(markers.MakeDefinition("ioc:autowired:constructor", markers.DescribesType, ""))   // constructor
	ConfigurationMarker = markers.Must(markers.MakeDefinition("ioc:autowired:configuration", markers.DescribesType, "")) // configuration -> @Configuration
	BeanMarker          = markers.Must(markers.MakeDefinition("ioc:autowired:bean", markers.DescribesType, ""))          // bean -> @Bean
)

type Generator struct {
	HeaderFile string `marker:",optional"`
	Year       string `marker:",optional"`
}

// CheckFilter genall.NeedsTypeChecking
func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		return IsNotInterface(node)
	}
}

func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into); err != nil {
		// TODO
		return err
	}
	return nil
}

func (d Generator) Generate(ctx *genall.GenerationContext) error {
	var headerText string

	if d.HeaderFile != "" {
		headerBytes, err := ctx.ReadFile(d.HeaderFile)
		if err != nil {
			return err
		}
		headerText = string(headerBytes)
	}

	iocCtx := NewIocGoGenerationContext(ctx, headerText)

	for _, root := range ctx.Roots {
		outContents := iocCtx.handlePackage(root)
		if outContents == nil {
			continue
		}

		WriteTo(ctx, root, outContents, "zz_generated.iocgo.go")
	}

	return nil
}

type IocGoGenerationContext struct {
	Collector  *markers.Collector
	Checker    *loader.TypeChecker
	HeaderText string
}

func (ctx *IocGoGenerationContext) handlePackage(root *loader.Package) []byte {
	typeInfos := make([]*markers.TypeInfo, 0)
	if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		typeInfos = append(typeInfos, info)
	}); err != nil {
		root.AddError(err)
		return nil
	}

	fire := false
	for _, info := range typeInfos {
		if len(info.Markers["ioc:autowired"]) != 0 {
			fire = true
			break
		}
	}
	if !fire {
		return nil
	}

	ctx.Checker.Check(root)
	root.NeedTypesInfo()

	imports := &Imports{
		byPath: make(map[string]string),
		byName: make(map[string]string),
		pkg:    root,
	}
	imports.byName[root.Name] = root.Name

	outContent := new(bytes.Buffer)
	bdm := NewBeanDefinitionMaker(root, imports, NewWriter(outContent), typeInfos)

	bdm.Generate()

	outBytes := outContent.Bytes()

	outContent = new(bytes.Buffer)
	writeHeader(root, outContent, root.Name, imports, ctx.HeaderText)
	writeMethods(root, outContent, outBytes)

	outBytes = outContent.Bytes()
	formattedBytes, err := format.Source(outBytes)
	if err != nil {
		root.AddError(err)
	} else {
		outBytes = formattedBytes
	}

	return outBytes
}

func NewIocGoGenerationContext(ctx *genall.GenerationContext, headerText string) IocGoGenerationContext {
	return IocGoGenerationContext{
		Collector:  ctx.Collector,
		Checker:    ctx.Checker,
		HeaderText: headerText,
	}
}

func IsInterface(node ast.Node) bool {
	_, ok := node.(*ast.InterfaceType)

	return ok
}

func IsNotInterface(node ast.Node) bool {
	return !IsInterface(node)
}

func writeHeader(pkg *loader.Package, out io.Writer, packageName string, imports *Imports, headerText string) {
	// TODO
}

func writeMethods(pkg *loader.Package, out io.Writer, outBuffer []byte) {
	// TODO
}
