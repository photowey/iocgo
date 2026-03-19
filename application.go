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
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/photowey/iocgo/events"
	"github.com/photowey/iocgo/internal/beanfactory"
)

type Application struct {
	factory beanfactory.BeanFactory
	events  events.Dispatcher
	mu      sync.Mutex
	booted  bool
}

// New creates a new application container with built-in event infrastructure.
//
// The container owns the default dispatcher so listener discovery, lifecycle
// events, and injected publishers all share the same root event mechanism.
func New() *Application {
	dispatcher := events.NewSyncDispatcher()
	factory := beanfactory.New()
	factory.SetEventPublisher(dispatcher)
	app := &Application{
		factory: factory,
		events:  dispatcher,
	}
	_ = app.factory.Register(
		Define[events.Publisher](events.PublisherBeanName, func(context.Context, Resolver) (events.Publisher, error) {
			return app.events, nil
		}),
		Define[events.Dispatcher](events.DispatcherBeanName, func(context.Context, Resolver) (events.Dispatcher, error) {
			return app.events, nil
		}),
	)
	return app
}

// Register adds bean definitions to the container before boot.
//
// Registration stays explicit so generated registrars, starters, and
// application code all contribute through the same deterministic path.
func (app *Application) Register(defs ...BeanDefinition) error {
	return app.factory.Register(defs...)
}

// Boot loads generated registrars, discovers listener beans, publishes
// container boot events, and eagerly creates non-lazy singletons.
//
// Boot is where the container turns metadata into a running dependency graph.
// Keeping this phase explicit makes startup ordering and failure modes easier to
// reason about than implicit lazy discovery.
func (app *Application) Boot(ctx context.Context) error {
	app.mu.Lock()
	if app.booted {
		app.mu.Unlock()
		return nil
	}
	app.mu.Unlock()

	bootstrap := DefaultBootstrapRegistry().Snapshot()
	for _, registrar := range bootstrap.BeanRegistrars {
		if err := registrar.Register(app); err != nil {
			return fmt.Errorf("register beans for module %q: %w", registrar.Module, err)
		}
	}
	for _, registrar := range bootstrap.StarterRegistrars {
		if err := registrar.Register(app); err != nil {
			return fmt.Errorf("register starter %q: %w", registrar.Module, err)
		}
	}
	if err := app.registerListenerBeans(ctx); err != nil {
		return err
	}
	if err := app.events.Publish(ctx, events.ContainerBootingEvent{At: time.Now()}); err != nil {
		return err
	}
	if err := app.factory.PreInstantiateSingletons(ctx); err != nil {
		return err
	}

	app.mu.Lock()
	app.booted = true
	app.mu.Unlock()
	return app.events.Publish(ctx, events.ContainerBootedEvent{At: time.Now()})
}

// Shutdown destroys managed singletons and publishes container shutdown events.
//
// Shutdown is explicit so resource cleanup and destruction order stay
// deterministic instead of depending on process exit timing.
func (app *Application) Shutdown(ctx context.Context) error {
	app.mu.Lock()
	if !app.booted {
		app.mu.Unlock()
		return nil
	}
	app.booted = false
	app.mu.Unlock()
	if err := app.events.Publish(ctx, events.ContainerShuttingDownEvent{At: time.Now()}); err != nil {
		return err
	}
	if err := app.factory.DestroySingletons(ctx); err != nil {
		return err
	}
	return app.events.Publish(ctx, events.ContainerShutdownEvent{At: time.Now()})
}

// Resolve resolves a single bean by exposed type and optional name.
func (app *Application) Resolve(ctx context.Context, typ reflect.Type, name string) (any, error) {
	return app.factory.Resolve(ctx, typ, name)
}

// ResolveAll resolves all beans exposing a given type using the shared ordering model.
func (app *Application) ResolveAll(ctx context.Context, typ reflect.Type) ([]any, error) {
	return app.factory.ResolveAll(ctx, typ)
}

// Definitions returns the current bean definition set in stable registration order.
func (app *Application) Definitions() []BeanDefinition {
	return app.factory.Definitions()
}

// Publish publishes an event through the built-in container dispatcher.
func (app *Application) Publish(ctx context.Context, event any) error {
	return app.events.Publish(ctx, event)
}

// Dispatcher exposes the built-in dispatcher for advanced framework integration.
func (app *Application) Dispatcher() events.Dispatcher {
	return app.events
}

func (app *Application) registerListenerBeans(ctx context.Context) error {
	listenerType := TypeOf[events.Listener]()
	for _, def := range app.factory.Definitions() {
		if !def.Exposes(listenerType) {
			continue
		}
		bean, err := app.Resolve(ctx, listenerType, def.Name)
		if err != nil {
			return fmt.Errorf("resolve listener bean %q: %w", def.Name, err)
		}
		listener, ok := bean.(events.Listener)
		if !ok {
			return fmt.Errorf("bean %q does not implement events.Listener", def.Name)
		}
		app.events.Register(listener)
	}
	return nil
}
