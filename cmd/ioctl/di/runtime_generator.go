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
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"sort"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type runtimeCodeBuilder struct {
	root           *loader.Package
	packageMarkers markers.MarkerValues
	typeInfos      []*markers.TypeInfo
	functions      []functionInfo
	headerText     string
	generatedAt    string
	usedImports    map[string]string
	fileImports    map[*ast.File]map[string]string
}

func newRuntimeCodeBuilder(root *loader.Package, packageMarkers markers.MarkerValues, typeInfos []*markers.TypeInfo, functions []functionInfo, headerText, generatedAt string) *runtimeCodeBuilder {
	builder := &runtimeCodeBuilder{
		root:           root,
		packageMarkers: packageMarkers,
		typeInfos:      typeInfos,
		functions:      functions,
		headerText:     headerText,
		generatedAt:    generatedAt,
		usedImports:    map[string]string{"iocgo": "github.com/photowey/iocgo", "context": "context"},
		fileImports:    make(map[*ast.File]map[string]string),
	}
	for _, file := range root.Syntax {
		builder.fileImports[file] = collectFileImports(file)
	}
	return builder
}

func (b *runtimeCodeBuilder) Build() ([]byte, error) {
	components := make([]string, 0)
	constructors := make(map[string]functionInfo)
	configurations := make(map[string]*markers.TypeInfo)
	methodsByReceiver := make(map[string][]functionInfo)
	freeBeans := make([]functionInfo, 0)

	for _, fn := range b.functions {
		if fn.Decl.Recv == nil {
			constructors[fn.Decl.Name.Name] = fn
		} else {
			receiver := receiverTypeName(fn.Decl)
			methodsByReceiver[receiver] = append(methodsByReceiver[receiver], fn)
		}
		if markerExists(fn.Markers, BeanMarker.Name) && fn.Decl.Recv == nil {
			freeBeans = append(freeBeans, fn)
		}
	}

	for _, info := range b.typeInfos {
		if err := b.validateTypeInfo(info); err != nil {
			return nil, err
		}
		if markerExists(info.Markers, ComponentMarker.Name) {
			definition, err := b.buildComponentDefinition(info, constructors)
			if err != nil {
				return nil, err
			}
			components = append(components, definition)
		}
		if markerExists(info.Markers, ConfigurationMarker.Name) {
			configurations[info.Name] = info
			definition, err := b.buildConfigurationDefinition(info)
			if err != nil {
				return nil, err
			}
			components = append(components, definition)
		}
	}

	for configName, info := range configurations {
		for _, fn := range methodsByReceiver[configName] {
			if !markerExists(fn.Markers, BeanMarker.Name) {
				continue
			}
			definition, err := b.buildBeanFunctionDefinition(fn, info)
			if err != nil {
				return nil, err
			}
			components = append(components, definition)
		}
	}

	for _, fn := range freeBeans {
		definition, err := b.buildBeanFunctionDefinition(fn, nil)
		if err != nil {
			return nil, err
		}
		components = append(components, definition)
	}

	var out bytes.Buffer
	b.writeHeader(&out)
	b.writeImports(&out)
	b.writeRegistrar(&out, components)
	b.writeInit(&out)
	return out.Bytes(), nil
}

func (b *runtimeCodeBuilder) validateTypeInfo(info *markers.TypeInfo) error {
	isComponent := markerExists(info.Markers, ComponentMarker.Name)
	isConfiguration := markerExists(info.Markers, ConfigurationMarker.Name)
	hasConfigurationProperties := markerExists(info.Markers, ConfigurationPropertiesMarker.Name) || bindingPrefix(info.Markers) != ""
	if markerExists(info.Markers, ConfigurationPropertiesMarker.Name) && bindingPrefix(info.Markers) == "" {
		return fmt.Errorf("type %s uses +ioc:configuration-properties without a prefix", info.Name)
	}
	if hasConfigurationProperties && !isComponent && !isConfiguration {
		return fmt.Errorf("type %s uses configuration-properties markers but is neither a component nor a configuration", info.Name)
	}
	if isConfiguration && !isStructType(info.RawSpec.Type) {
		return fmt.Errorf("configuration type %s must be a struct", info.Name)
	}
	if hasConfigurationProperties && !isStructType(info.RawSpec.Type) {
		return fmt.Errorf("configuration-properties type %s must be a struct", info.Name)
	}
	if isComponent && !isStructType(info.RawSpec.Type) && markerString(info.Markers, ConstructorMarker.Name) == "" {
		return fmt.Errorf("component type %s must be a struct unless a constructor is specified", info.Name)
	}
	return nil
}

func isStructType(expr ast.Expr) bool {
	_, ok := expr.(*ast.StructType)
	return ok
}

func (b *runtimeCodeBuilder) writeHeader(out *bytes.Buffer) {
	out.WriteString("// Code generated by iocgo ctl (ioctl). DO NOT EDIT.\n")
	out.WriteString("//\n")
	if _, err := fmt.Fprintf(out, "// Generated at: %s\n", b.generatedAt); err != nil {
		panic(err)
	}
	if _, err := fmt.Fprintf(out, "// Source package: %s\n", b.root.PkgPath); err != nil {
		panic(err)
	}
	out.WriteString("//\n")
	if strings.TrimSpace(b.headerText) != "" {
		out.WriteString(strings.TrimRight(b.headerText, "\n"))
		out.WriteString("\n\n")
	} else {
		out.WriteString("\n")
	}
	out.WriteString("package ")
	out.WriteString(b.root.Name)
	out.WriteString("\n\n")
}

func (b *runtimeCodeBuilder) writeImports(out *bytes.Buffer) {
	aliases := make([]string, 0, len(b.usedImports))
	for alias := range b.usedImports {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	out.WriteString("import (\n")
	for _, alias := range aliases {
		path := b.usedImports[alias]
		if alias == importDefaultAlias(path) {
			if _, err := fmt.Fprintf(out, "\t%q\n", path); err != nil {
				panic(err)
			}
			continue
		}
		if _, err := fmt.Fprintf(out, "\t%s %q\n", alias, path); err != nil {
			panic(err)
		}
	}
	out.WriteString(")\n\n")
}

func (b *runtimeCodeBuilder) writeRegistrar(out *bytes.Buffer, definitions []string) {
	out.WriteString("func RegisterIocgoBeans(reg iocgo.Registry) error {\n")
	if len(definitions) == 0 {
		out.WriteString("\treturn nil\n")
		out.WriteString("}\n\n")
		return
	}
	out.WriteString("\treturn reg.Register(\n")
	for _, definition := range definitions {
		out.WriteString(indent(definition, "\t\t"))
		out.WriteString(",\n")
	}
	out.WriteString("\t)\n")
	out.WriteString("}\n\n")
}

func (b *runtimeCodeBuilder) writeInit(out *bytes.Buffer) {
	registerBeans := markerExists(b.packageMarkers, InitRegisterMarker.Name)
	registerStarter := markerExists(b.packageMarkers, StarterMarker.Name)
	if !registerBeans && !registerStarter {
		return
	}
	out.WriteString("func init() {\n")
	if registerBeans {
		if _, err := fmt.Fprintf(out, "\tiocgo.RegisterBeans(%q, RegisterIocgoBeans)\n", b.root.PkgPath); err != nil {
			panic(err)
		}
	}
	if registerStarter {
		if _, err := fmt.Fprintf(out, "\tiocgo.RegisterStarter(%q, RegisterIocgoBeans)\n", b.root.PkgPath); err != nil {
			panic(err)
		}
	}
	out.WriteString("}\n")
}

func (b *runtimeCodeBuilder) buildComponentDefinition(info *markers.TypeInfo, constructors map[string]functionInfo) (string, error) {
	beanName := markerString(info.Markers, ComponentNameMarker.Name)
	if beanName == "" {
		beanName = lowerCamel(info.Name)
	}
	typeExpr := "*" + info.Name
	factoryBody := &strings.Builder{}
	constructorName := markerString(info.Markers, ConstructorMarker.Name)
	if constructorName != "" {
		constructor, ok := constructors[constructorName]
		if !ok {
			return "", fmt.Errorf("component %s references unknown constructor %s", info.Name, constructorName)
		}
		call, err := b.buildFunctionInvocation(constructor, nil)
		if err != nil {
			return "", err
		}
		factoryBody.WriteString(call)
	} else {
		if _, err := fmt.Fprintf(factoryBody, "bean := &%s{}\nreturn bean, nil", info.Name); err != nil {
			panic(err)
		}
	}
	factoryBodyString := b.wrapBindingFactoryBody(beanName, "*"+info.Name, info.Markers, factoryBody.String())
	injector, err := b.buildFieldInjector(info)
	if err != nil {
		return "", err
	}
	return b.renderDefinition(typeExpr, beanName, info.Markers, info.RawFile, info.Name, factoryBodyString, injector), nil
}

func (b *runtimeCodeBuilder) buildConfigurationDefinition(info *markers.TypeInfo) (string, error) {
	beanName := lowerCamel(info.Name)
	typeExpr := "*" + info.Name
	injector, err := b.buildFieldInjector(info)
	if err != nil {
		return "", err
	}
	factoryBody := b.wrapBindingFactoryBody(beanName, typeExpr, info.Markers, fmt.Sprintf("bean := &%s{}\nreturn bean, nil", info.Name))
	return b.renderDefinition(typeExpr, beanName, info.Markers, info.RawFile, info.Name, factoryBody, injector), nil
}

func (b *runtimeCodeBuilder) buildBeanFunctionDefinition(fn functionInfo, config *markers.TypeInfo) (string, error) {
	beanName := markerString(fn.Markers, BeanNameMarker.Name)
	if beanName == "" {
		beanName = lowerCamel(fn.Decl.Name.Name)
	}
	returnType, hasError, err := b.parseReturnSignature(fn)
	if err != nil {
		return "", err
	}
	call, err := b.buildFunctionInvocation(fn, config)
	if err != nil {
		return "", err
	}
	_ = hasError
	return b.renderDefinition(returnType, beanName, fn.Markers, fn.File, fn.Decl.Name.Name, call, ""), nil
}

func (b *runtimeCodeBuilder) renderDefinition(typeExpr, beanName string, markerValues markers.MarkerValues, file *ast.File, symbol, factoryBody, injector string) string {
	options := make([]string, 0)
	scopeValue := markerString(markerValues, ScopeMarker.Name)
	if scopeValue == "prototype" {
		options = append(options, "iocgo.WithScope(iocgo.Prototype)")
	} else {
		options = append(options, "iocgo.WithScope(iocgo.Singleton)")
	}
	if markerExists(markerValues, PrimaryMarker.Name) {
		options = append(options, "iocgo.WithPrimary()")
	}
	if markerExists(markerValues, LazyMarker.Name) {
		options = append(options, "iocgo.WithLazy()")
	}
	for _, expose := range parseList(markerString(markerValues, ExposeMarker.Name)) {
		b.trackTypeReference(file, expose)
		options = append(options, fmt.Sprintf("iocgo.WithExposedTypes(iocgo.TypeOf[%s]())", expose))
	}
	if injector != "" {
		options = append(options, injector)
	}
	if initHook := markerString(markerValues, InitMarker.Name); initHook != "" {
		options = append(options, fmt.Sprintf("iocgo.WithInitHook(func(ctx context.Context, bean any) error {\n\ttarget := bean.(%s)\n\ttarget.%s()\n\treturn nil\n})", typeExpr, initHook))
	}
	if destroyHook := markerString(markerValues, DestroyMarker.Name); destroyHook != "" {
		options = append(options, fmt.Sprintf("iocgo.WithDestroyHook(func(ctx context.Context, bean any) error {\n\ttarget := bean.(%s)\n\ttarget.%s()\n\treturn nil\n})", typeExpr, destroyHook))
	}
	fileName := b.root.Fset.Position(file.Package).Filename
	options = append(options, fmt.Sprintf("iocgo.WithSource(iocgo.SourceInfo{Package: %q, File: %q, Symbol: %q})", b.root.PkgPath, fileName, symbol))

	definition := &strings.Builder{}
	if _, err := fmt.Fprintf(definition, "iocgo.Define[%s](%q, func(ctx context.Context, resolver iocgo.Resolver) (%s, error) {\n", typeExpr, beanName, typeExpr); err != nil {
		panic(err)
	}
	definition.WriteString(indent(factoryBody, "\t"))
	definition.WriteString("\n}")
	for _, option := range options {
		definition.WriteString(",\n")
		definition.WriteString(indent(option, "\t"))
	}
	definition.WriteString(")")
	return definition.String()
}

func (b *runtimeCodeBuilder) wrapBindingFactoryBody(beanName, typeExpr string, markerValues markers.MarkerValues, factoryBody string) string {
	prefix := bindingPrefix(markerValues)
	if prefix == "" {
		return factoryBody
	}
	b.usedImports["starterNemo"] = "github.com/photowey/iocgo/pkg/starters/nemo"

	binding := strings.Join([]string{
		"binder, err := iocgo.Get[*starterNemo.Binder](ctx, resolver, starterNemo.BinderBeanName)",
		"if err != nil {",
		fmt.Sprintf("\treturn nil, starterNemo.WrapBindError(%q, %q, %q, err)", beanName, typeExpr, prefix),
		"}",
		fmt.Sprintf("if err := binder.Bind(%q, bean); err != nil {", prefix),
		fmt.Sprintf("\treturn nil, starterNemo.WrapBindError(%q, %q, %q, err)", beanName, typeExpr, prefix),
		"}",
	}, "\n")

	return strings.Replace(factoryBody, "return bean, nil", binding+"\nreturn bean, nil", 1)
}

func bindingPrefix(markerValues markers.MarkerValues) string {
	if prefix := markerString(markerValues, ConfigurationPropertiesPrefixMarker.Name); prefix != "" {
		return prefix
	}
	return markerString(markerValues, BindPrefixMarker.Name)
}

func (b *runtimeCodeBuilder) buildFieldInjector(info *markers.TypeInfo) (string, error) {
	lines := make([]string, 0)
	for _, field := range info.Fields {
		if !markerExists(field.Markers, InjectMarker.Name) && !markerExists(field.Markers, InjectNameMarker.Name) {
			continue
		}
		if field.Name == "" {
			continue
		}
		typeExpr := b.exprString(field.RawField.Type)
		b.trackExprImports(info.RawFile, field.RawField.Type)
		beanName := markerString(field.Markers, InjectNameMarker.Name)
		varName := lowerCamel(field.Name) + "Dependency"
		if beanName == "" {
			lines = append(lines, fmt.Sprintf("%s, err := iocgo.Get[%s](ctx, resolver)", varName, typeExpr))
		} else {
			lines = append(lines, fmt.Sprintf("%s, err := iocgo.Get[%s](ctx, resolver, %q)", varName, typeExpr, beanName))
		}
		lines = append(lines, "if err != nil {", "\treturn err", "}")
		lines = append(lines, fmt.Sprintf("target.%s = %s", field.Name, varName))
	}
	if len(lines) == 0 {
		return "", nil
	}
	builder := &strings.Builder{}
	if _, err := fmt.Fprintf(builder, "iocgo.WithInjector(func(ctx context.Context, resolver iocgo.Resolver, bean any) error {\n\ttarget := bean.(*%s)\n", info.Name); err != nil {
		panic(err)
	}
	for _, line := range lines {
		builder.WriteString("\t")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("\treturn nil\n})")
	return builder.String(), nil
}

func (b *runtimeCodeBuilder) buildFunctionInvocation(fn functionInfo, config *markers.TypeInfo) (string, error) {
	args, err := b.buildCallArgs(fn)
	if err != nil {
		return "", err
	}
	callTarget := fn.Decl.Name.Name
	if config != nil {
		configBeanName := lowerCamel(config.Name)
		configType := "*" + config.Name
		args = append([]string{fmt.Sprintf("cfg, err := iocgo.Get[%s](ctx, resolver, %q)", configType, configBeanName), "if err != nil {", "\treturn nil, err", "}"}, args...)
		callTarget = "cfg." + fn.Decl.Name.Name
	}
	resultType, hasError, err := b.parseReturnSignature(fn)
	if err != nil {
		return "", err
	}
	argNames := callArgumentNames(fn.Decl)
	builder := &strings.Builder{}
	for _, line := range args {
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	if hasError {
		fmt.Fprintf(builder, "bean, err := %s(%s)\n", callTarget, strings.Join(argNames, ", "))
		builder.WriteString("if err != nil {\n\treturn nil, err\n}\n")
		builder.WriteString("return bean, nil")
		return builder.String(), nil
	}
	fmt.Fprintf(builder, "bean := %s(%s)\n", callTarget, strings.Join(argNames, ", "))
	builder.WriteString("return bean, nil")
	_ = resultType
	return builder.String(), nil
}

func (b *runtimeCodeBuilder) buildCallArgs(fn functionInfo) ([]string, error) {
	bindings := parseNamedBindings(markerString(fn.Markers, InjectNamedMarker.Name))
	lines := make([]string, 0)
	params := fn.Decl.Type.Params
	if params == nil {
		return lines, nil
	}
	for _, field := range params.List {
		typeExpr := b.exprString(field.Type)
		for _, name := range parameterNames(field) {
			if typeExpr == "context.Context" {
				continue
			}
			b.trackExprImports(fn.File, field.Type)
			if beanName := bindings[name]; beanName != "" {
				lines = append(lines, fmt.Sprintf("%s, err := iocgo.Get[%s](ctx, resolver, %q)", name, typeExpr, beanName))
			} else {
				lines = append(lines, fmt.Sprintf("%s, err := iocgo.Get[%s](ctx, resolver)", name, typeExpr))
			}
			lines = append(lines, "if err != nil {", "\treturn nil, err", "}")
		}
	}
	return lines, nil
}

func (b *runtimeCodeBuilder) parseReturnSignature(fn functionInfo) (string, bool, error) {
	results := fn.Decl.Type.Results
	if results == nil || len(results.List) == 0 {
		return "", false, fmt.Errorf("bean function %s must return a bean", fn.Decl.Name.Name)
	}
	if len(results.List) == 1 {
		typeExpr := b.exprString(results.List[0].Type)
		b.trackExprImports(fn.File, results.List[0].Type)
		return typeExpr, false, nil
	}
	if len(results.List) == 2 && b.exprString(results.List[1].Type) == "error" {
		typeExpr := b.exprString(results.List[0].Type)
		b.trackExprImports(fn.File, results.List[0].Type)
		return typeExpr, true, nil
	}
	return "", false, fmt.Errorf("bean function %s must return T or (T, error)", fn.Decl.Name.Name)
}

func (b *runtimeCodeBuilder) exprString(expr ast.Expr) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, token.NewFileSet(), expr)
	return buf.String()
}

func (b *runtimeCodeBuilder) trackExprImports(file *ast.File, expr ast.Expr) {
	imports := b.fileImports[file]
	ast.Inspect(expr, func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}
		path, ok := imports[ident.Name]
		if ok {
			b.usedImports[ident.Name] = path
		}
		return true
	})
}

func (b *runtimeCodeBuilder) trackTypeReference(file *ast.File, typeExpr string) {
	if !strings.Contains(typeExpr, ".") {
		return
	}
	alias := strings.Split(typeExpr, ".")[0]
	if path, ok := b.fileImports[file][alias]; ok {
		b.usedImports[alias] = path
	}
}

func collectFileImports(file *ast.File) map[string]string {
	imports := make(map[string]string)
	for _, spec := range file.Imports {
		path := strings.Trim(spec.Path.Value, "\"")
		alias := importDefaultAlias(path)
		if spec.Name != nil && spec.Name.Name != "" {
			alias = spec.Name.Name
		}
		imports[alias] = path
	}
	return imports
}

func importDefaultAlias(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func receiverTypeName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}
	switch expr := fn.Recv.List[0].Type.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		if ident, ok := expr.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

func lowerCamel(src string) string {
	if src == "" {
		return src
	}
	return strings.ToLower(src[:1]) + src[1:]
}

func parseList(src string) []string {
	src = strings.TrimSpace(src)
	src = strings.TrimPrefix(src, "{")
	src = strings.TrimSuffix(src, "}")
	if src == "" {
		return nil
	}
	parts := strings.Split(src, ";")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}

func parseNamedBindings(src string) map[string]string {
	bindings := make(map[string]string)
	for _, part := range parseList(src) {
		items := strings.SplitN(part, ":", 2)
		if len(items) != 2 {
			continue
		}
		bindings[strings.TrimSpace(items[0])] = strings.TrimSpace(items[1])
	}
	return bindings
}

func parameterNames(field *ast.Field) []string {
	if len(field.Names) == 0 {
		return []string{"arg"}
	}
	names := make([]string, 0, len(field.Names))
	for _, ident := range field.Names {
		names = append(names, ident.Name)
	}
	return names
}

func callArgumentNames(fn *ast.FuncDecl) []string {
	args := make([]string, 0)
	if fn.Type.Params == nil {
		return args
	}
	for _, field := range fn.Type.Params.List {
		typeExpr := exprToString(field.Type)
		for _, name := range parameterNames(field) {
			if typeExpr == "context.Context" {
				args = append(args, "ctx")
				continue
			}
			args = append(args, name)
		}
	}
	return args
}

func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, token.NewFileSet(), expr)
	return buf.String()
}

func indent(src, prefix string) string {
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}
