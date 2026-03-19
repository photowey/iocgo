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

package beandefinition

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/photowey/iocgo/internal/scope"
)

type Resolver interface {
	Resolve(ctx context.Context, typ reflect.Type, name string) (any, error)
	ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error)
}

type FactoryFunc func(ctx context.Context, resolver Resolver) (any, error)
type InjectorFunc func(ctx context.Context, resolver Resolver, bean any) error
type Hook func(ctx context.Context, bean any) error

type DependencyKind uint8

const (
	DependencyDirect DependencyKind = iota + 1
	DependencyProvider
	DependencyLazy
)

type Dependency struct {
	Name     string
	Type     reflect.Type
	BeanName string
	Optional bool
	Kind     DependencyKind
}

type SourceInfo struct {
	Package string
	File    string
	Symbol  string
}

// BeanDefinition describes how a managed bean should be created and resolved.
type BeanDefinition struct {
	Name         string
	Scope        scope.Scope
	Primary      bool
	Lazy         bool
	ExposedTypes []reflect.Type
	Dependencies []Dependency
	Factory      FactoryFunc
	Injector     InjectorFunc
	InitHook     Hook
	DestroyHook  Hook
	Source       SourceInfo
}

func (bd BeanDefinition) Validate() error {
	if bd.Name == "" {
		return fmt.Errorf("bean definition requires a name")
	}
	if !bd.Scope.Valid() {
		return fmt.Errorf("bean %q uses invalid scope %q", bd.Name, bd.Scope.String())
	}
	if len(bd.ExposedTypes) == 0 {
		return fmt.Errorf("bean %q must expose at least one type", bd.Name)
	}
	if bd.Factory == nil {
		return fmt.Errorf("bean %q requires a factory", bd.Name)
	}
	for _, typ := range bd.ExposedTypes {
		if typ == nil {
			return fmt.Errorf("bean %q contains a nil exposed type", bd.Name)
		}
	}
	return nil
}

func (bd BeanDefinition) Exposes(typ reflect.Type) bool {
	if typ == nil {
		return false
	}
	for _, candidate := range bd.ExposedTypes {
		if candidate == typ {
			return true
		}
	}
	return false
}

func (bd BeanDefinition) Description() string {
	parts := []string{bd.Name, bd.Scope.String()}
	if bd.Source.Package != "" || bd.Source.Symbol != "" {
		parts = append(parts, fmt.Sprintf("%s:%s", bd.Source.Package, bd.Source.Symbol))
	}
	return strings.Join(parts, "|")
}
