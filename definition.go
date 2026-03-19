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

package iocgo

import (
	"context"
	"reflect"

	"github.com/photowey/iocgo/internal/beandefinition"
	"github.com/photowey/iocgo/internal/scope"
)

type Scope = scope.Scope

const (
	Singleton = scope.Singleton
	Prototype = scope.Prototype
)

type BeanDefinition = beandefinition.BeanDefinition
type Dependency = beandefinition.Dependency
type DependencyKind = beandefinition.DependencyKind
type SourceInfo = beandefinition.SourceInfo
type FactoryFunc = beandefinition.FactoryFunc
type InjectorFunc = beandefinition.InjectorFunc
type Hook = beandefinition.Hook
type Resolver = beandefinition.Resolver

type Registry interface {
	Register(defs ...BeanDefinition) error
}

type TypeResolver interface {
	Resolve(ctx context.Context, typ reflect.Type, name string) (any, error)
	ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error)
}

type RegistrarFunc func(reg Registry) error

type DefinitionOption func(def *BeanDefinition)

// WithScope customizes bean scope while keeping scope selection explicit in definition metadata.
func WithScope(value Scope) DefinitionOption {
	return func(def *BeanDefinition) {
		def.Scope = value
	}
}

// WithPrimary marks a bean as the preferred candidate when type resolution is ambiguous.
func WithPrimary() DefinitionOption {
	return func(def *BeanDefinition) {
		def.Primary = true
	}
}

// WithLazy keeps singleton creation deferred until first resolution.
func WithLazy() DefinitionOption {
	return func(def *BeanDefinition) {
		def.Lazy = true
	}
}

// WithExposedTypes declares the full type surface a bean should be resolved as.
func WithExposedTypes(types ...reflect.Type) DefinitionOption {
	return func(def *BeanDefinition) {
		def.ExposedTypes = append(def.ExposedTypes, types...)
	}
}

func WithDependencies(deps ...Dependency) DefinitionOption {
	return func(def *BeanDefinition) {
		def.Dependencies = append(def.Dependencies, deps...)
	}
}

// WithInjector attaches post-construction dependency injection logic to a bean definition.
func WithInjector(injector InjectorFunc) DefinitionOption {
	return func(def *BeanDefinition) {
		def.Injector = injector
	}
}

// WithInitHook attaches an explicit initialization hook to a bean definition.
func WithInitHook(hook Hook) DefinitionOption {
	return func(def *BeanDefinition) {
		def.InitHook = hook
	}
}

// WithDestroyHook attaches an explicit destruction hook to a bean definition.
func WithDestroyHook(hook Hook) DefinitionOption {
	return func(def *BeanDefinition) {
		def.DestroyHook = hook
	}
}

func WithSource(source SourceInfo) DefinitionOption {
	return func(def *BeanDefinition) {
		def.Source = source
	}
}

// Define creates a bean definition from a typed factory function.
//
// Define is the main bridge between generated code and the runtime container. It
// captures enough typed information up front that the runtime can stay concrete
// and deterministic instead of reflecting over arbitrary factories at boot time.
func Define[T any](name string, factory func(context.Context, Resolver) (T, error), opts ...DefinitionOption) BeanDefinition {
	def := BeanDefinition{
		Name:         name,
		Scope:        Singleton,
		ExposedTypes: []reflect.Type{TypeOf[T]()},
		Factory: func(ctx context.Context, resolver Resolver) (any, error) {
			return factory(ctx, resolver)
		},
	}
	for _, opt := range opts {
		opt(&def)
	}
	def.ExposedTypes = uniqueTypes(def.ExposedTypes)
	return def
}

func uniqueTypes(types []reflect.Type) []reflect.Type {
	seen := make(map[reflect.Type]struct{}, len(types))
	unique := make([]reflect.Type, 0, len(types))
	for _, typ := range types {
		if typ == nil {
			continue
		}
		if _, ok := seen[typ]; ok {
			continue
		}
		seen[typ] = struct{}{}
		unique = append(unique, typ)
	}
	return unique
}
