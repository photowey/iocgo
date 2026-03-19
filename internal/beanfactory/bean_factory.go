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

package beanfactory

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/photowey/iocgo/events"
	"github.com/photowey/iocgo/internal/beandefinition"
	"github.com/photowey/iocgo/internal/lifecycle"
	"github.com/photowey/iocgo/internal/scope"
	"github.com/photowey/iocgo/ordering"
)

type BeanFactory interface {
	Register(defs ...beandefinition.BeanDefinition) error
	Resolve(ctx context.Context, typ reflect.Type, name string) (any, error)
	ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error)
	PreInstantiateSingletons(ctx context.Context) error
	DestroySingletons(ctx context.Context) error
	Definitions() []beandefinition.BeanDefinition
	SetEventPublisher(events.Publisher)
}

type defaultFactory struct {
	mu                sync.RWMutex
	defsByName        map[string]beandefinition.BeanDefinition
	typeIndex         map[reflect.Type][]string
	registrationOrder []string
	singletons        map[string]any
	creationOrder     []string
	beanLocks         map[string]*sync.Mutex
	eventPublisher    events.Publisher
}

func New() BeanFactory {
	return &defaultFactory{
		defsByName: make(map[string]beandefinition.BeanDefinition),
		typeIndex:  make(map[reflect.Type][]string),
		singletons: make(map[string]any),
		beanLocks:  make(map[string]*sync.Mutex),
	}
}

func (f *defaultFactory) Register(defs ...beandefinition.BeanDefinition) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, def := range defs {
		if err := def.Validate(); err != nil {
			return err
		}
		if _, exists := f.defsByName[def.Name]; exists {
			return fmt.Errorf("duplicate bean definition %q", def.Name)
		}
		f.defsByName[def.Name] = def
		f.registrationOrder = append(f.registrationOrder, def.Name)
		for _, typ := range def.ExposedTypes {
			f.typeIndex[typ] = appendIfMissing(f.typeIndex[typ], def.Name)
		}
		if _, ok := f.beanLocks[def.Name]; !ok {
			f.beanLocks[def.Name] = &sync.Mutex{}
		}
	}

	return nil
}

func (f *defaultFactory) Resolve(ctx context.Context, typ reflect.Type, name string) (any, error) {
	def, err := f.selectDefinition(typ, name)
	if err != nil {
		return nil, err
	}
	return f.resolveDefinition(ctx, def.Name, nil)
}

func (f *defaultFactory) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	f.mu.RLock()
	names := append([]string(nil), f.typeIndex[typ]...)
	f.mu.RUnlock()
	if len(names) == 0 {
		return nil, fmt.Errorf("no beans found for type %s", typ)
	}

	beans := make([]any, 0, len(names))
	for _, name := range names {
		bean, err := f.resolveDefinition(ctx, name, nil)
		if err != nil {
			return nil, err
		}
		beans = append(beans, bean)
	}
	ordering.SortAny(beans)
	return beans, nil
}

func (f *defaultFactory) PreInstantiateSingletons(ctx context.Context) error {
	f.mu.RLock()
	names := append([]string(nil), f.registrationOrder...)
	f.mu.RUnlock()

	for _, name := range names {
		f.mu.RLock()
		def := f.defsByName[name]
		f.mu.RUnlock()
		if def.Scope == scope.Singleton && !def.Lazy {
			if _, err := f.resolveDefinition(ctx, name, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *defaultFactory) DestroySingletons(ctx context.Context) error {
	f.mu.RLock()
	order := append([]string(nil), f.creationOrder...)
	beans := make(map[string]any, len(f.singletons))
	for name, bean := range f.singletons {
		beans[name] = bean
	}
	defs := make(map[string]beandefinition.BeanDefinition, len(f.defsByName))
	for name, def := range f.defsByName {
		defs[name] = def
	}
	f.mu.RUnlock()

	var errs []string
	for i := len(order) - 1; i >= 0; i-- {
		name := order[i]
		bean, ok := beans[name]
		if !ok {
			continue
		}
		if disposable, ok := bean.(lifecycle.DisposableBean); ok {
			disposable.Destroy()
		}
		if hook := defs[name].DestroyHook; hook != nil {
			if err := hook(ctx, bean); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if publisher := f.publisher(); publisher != nil {
			if err := publisher.Publish(ctx, events.BeanDestroyedEvent{
				At:       time.Now(),
				BeanName: name,
				BeanType: reflect.TypeOf(bean),
			}); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	f.mu.Lock()
	f.singletons = make(map[string]any)
	f.creationOrder = nil
	f.mu.Unlock()

	if len(errs) > 0 {
		return fmt.Errorf("destroy singletons: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (f *defaultFactory) Definitions() []beandefinition.BeanDefinition {
	f.mu.RLock()
	defer f.mu.RUnlock()

	defs := make([]beandefinition.BeanDefinition, 0, len(f.registrationOrder))
	for _, name := range f.registrationOrder {
		defs = append(defs, f.defsByName[name])
	}
	return defs
}

func (f *defaultFactory) SetEventPublisher(publisher events.Publisher) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.eventPublisher = publisher
}

func (f *defaultFactory) selectDefinition(typ reflect.Type, name string) (beandefinition.BeanDefinition, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if name != "" {
		def, ok := f.defsByName[name]
		if !ok {
			return beandefinition.BeanDefinition{}, fmt.Errorf("bean %q is not registered", name)
		}
		if typ != nil && !def.Exposes(typ) {
			return beandefinition.BeanDefinition{}, fmt.Errorf("bean %q does not expose type %s", name, typ)
		}
		return def, nil
	}

	candidates := f.typeIndex[typ]
	if len(candidates) == 0 {
		return beandefinition.BeanDefinition{}, fmt.Errorf("no beans found for type %s", typ)
	}
	if len(candidates) == 1 {
		return f.defsByName[candidates[0]], nil
	}

	var primary *beandefinition.BeanDefinition
	for _, candidate := range candidates {
		def := f.defsByName[candidate]
		if !def.Primary {
			continue
		}
		if primary != nil {
			return beandefinition.BeanDefinition{}, fmt.Errorf("multiple primary beans found for type %s", typ)
		}
		clone := def
		primary = &clone
	}
	if primary != nil {
		return *primary, nil
	}

	return beandefinition.BeanDefinition{}, fmt.Errorf("multiple beans found for type %s: %s", typ, strings.Join(candidates, ", "))
}

func (f *defaultFactory) resolveDefinition(ctx context.Context, name string, trail []string) (any, error) {
	if contains(trail, name) {
		cycle := append(append([]string(nil), trail...), name)
		return nil, fmt.Errorf("bean cycle detected: %s", strings.Join(cycle, " -> "))
	}

	f.mu.RLock()
	def, ok := f.defsByName[name]
	cached, isCached := f.singletons[name]
	f.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("bean %q is not registered", name)
	}
	if isCached {
		return cached, nil
	}

	if def.Scope == scope.Singleton {
		lock := f.beanLock(name)
		lock.Lock()
		defer lock.Unlock()

		f.mu.RLock()
		cached, isCached = f.singletons[name]
		f.mu.RUnlock()
		if isCached {
			return cached, nil
		}
	}

	resolver := &factoryResolver{
		factory: f,
		trail:   append(append([]string(nil), trail...), name),
	}

	bean, err := def.Factory(ctx, resolver)
	if err != nil {
		return nil, fmt.Errorf("create bean %q: %w", name, err)
	}
	if def.Injector != nil {
		if err := def.Injector(ctx, resolver, bean); err != nil {
			return nil, fmt.Errorf("inject bean %q: %w", name, err)
		}
	}
	if initializing, ok := bean.(lifecycle.InitializingBean); ok {
		initializing.AfterPropertiesSet()
	}
	if def.InitHook != nil {
		if err := def.InitHook(ctx, bean); err != nil {
			return nil, fmt.Errorf("initialize bean %q: %w", name, err)
		}
	}

	if def.Scope == scope.Singleton {
		f.mu.Lock()
		if _, exists := f.singletons[name]; !exists {
			f.singletons[name] = bean
			f.creationOrder = append(f.creationOrder, name)
		}
		f.mu.Unlock()
	}

	if publisher := f.publisher(); publisher != nil {
		_ = publisher.Publish(ctx, events.BeanInitializedEvent{
			At:       time.Now(),
			BeanName: name,
			BeanType: reflect.TypeOf(bean),
		})
	}

	return bean, nil
}

func (f *defaultFactory) beanLock(name string) *sync.Mutex {
	f.mu.Lock()
	defer f.mu.Unlock()
	lock, ok := f.beanLocks[name]
	if !ok {
		lock = &sync.Mutex{}
		f.beanLocks[name] = lock
	}
	return lock
}

func (f *defaultFactory) publisher() events.Publisher {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.eventPublisher
}

type factoryResolver struct {
	factory *defaultFactory
	trail   []string
}

func (r *factoryResolver) Resolve(ctx context.Context, typ reflect.Type, name string) (any, error) {
	def, err := r.factory.selectDefinition(typ, name)
	if err != nil {
		return nil, err
	}
	return r.factory.resolveDefinition(ctx, def.Name, r.trail)
}

func (r *factoryResolver) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	return r.factory.ResolveAll(ctx, typ)
}

func appendIfMissing(src []string, value string) []string {
	for _, current := range src {
		if current == value {
			return src
		}
	}
	return append(src, value)
}

func contains(src []string, value string) bool {
	for _, current := range src {
		if current == value {
			return true
		}
	}
	return false
}
