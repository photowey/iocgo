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

package events

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/photowey/iocgo/ordering"
)

const (
	// PublisherBeanName is the built-in bean name used to inject the default
	// event publisher.
	PublisherBeanName = "iocgoEventPublisher"
	// DispatcherBeanName is the built-in bean name used to inject the default
	// event dispatcher.
	DispatcherBeanName = "iocgoEventDispatcher"
)

// Publisher publishes an event into the container event system.
type Publisher interface {
	Publish(context.Context, any) error
}

// Dispatcher extends Publisher with listener registration operations.
type Dispatcher interface {
	Publisher
	Register(Listener)
	RegisterAll(...Listener)
}

type Ordered = ordering.Ordered
type PriorityOrdered = ordering.PriorityOrdered

// Listener receives events of a specific runtime type.
type Listener interface {
	EventType() reflect.Type
	Handle(context.Context, any) error
}

// ListenerFunc adapts a function into a Listener.
type ListenerFunc struct {
	Type      reflect.Type
	Fn        func(context.Context, any) error
	SortOrder int
}

// PriorityListenerFunc adapts a function into a high-priority listener.
type PriorityListenerFunc struct {
	ListenerFunc
}

// EventType reports the runtime event type supported by this listener.
func (l ListenerFunc) EventType() reflect.Type {
	return l.Type
}

// Handle invokes the underlying listener function.
func (l ListenerFunc) Handle(ctx context.Context, event any) error {
	return l.Fn(ctx, event)
}

// Order returns the configured order for this listener function.
func (l ListenerFunc) Order() int {
	return l.SortOrder
}

// PriorityOrder marks this listener as part of the PriorityOrdered tier.
func (PriorityListenerFunc) PriorityOrder() {}

// NewListener creates a typed listener for a specific event payload type.
func NewListener[T any](fn func(context.Context, T) error, order ...int) Listener {
	sortOrder := 0
	if len(order) > 0 {
		sortOrder = order[0]
	}
	return ListenerFunc{
		Type: reflect.TypeOf((*T)(nil)).Elem(),
		Fn: func(ctx context.Context, event any) error {
			return fn(ctx, event.(T))
		},
		SortOrder: sortOrder,
	}
}

// NewPriorityListener creates a typed listener that participates in the
// PriorityOrdered tier.
func NewPriorityListener[T any](fn func(context.Context, T) error, order ...int) Listener {
	return PriorityListenerFunc{ListenerFunc: NewListener(fn, order...).(ListenerFunc)}
}

type SyncDispatcher struct {
	mu        sync.RWMutex
	seq       uint64
	listeners map[reflect.Type][]registeredListener
}

type registeredListener struct {
	listener Listener
	seq      uint64
}

// NewSyncDispatcher creates the default synchronous dispatcher implementation.
func NewSyncDispatcher(listeners ...Listener) *SyncDispatcher {
	dispatcher := &SyncDispatcher{
		listeners: make(map[reflect.Type][]registeredListener),
	}
	dispatcher.RegisterAll(listeners...)
	return dispatcher
}

// Register registers a listener with the dispatcher.
func (d *SyncDispatcher) Register(listener Listener) {
	if listener == nil || listener.EventType() == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seq++
	eventType := listener.EventType()
	d.listeners[eventType] = append(d.listeners[eventType], registeredListener{
		listener: listener,
		seq:      d.seq,
	})
}

// RegisterAll registers multiple listeners with the dispatcher.
func (d *SyncDispatcher) RegisterAll(listeners ...Listener) {
	for _, listener := range listeners {
		d.Register(listener)
	}
}

// Publish dispatches an event synchronously to matching listeners using the
// shared PriorityOrdered / Ordered model.
func (d *SyncDispatcher) Publish(ctx context.Context, event any) error {
	if event == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	eventType := reflect.TypeOf(event)
	d.mu.RLock()
	current := append([]registeredListener(nil), d.listeners[eventType]...)
	d.mu.RUnlock()
	// Preserve registration order for ties after ordering tiers and order values.
	sort.SliceStable(current, func(i, j int) bool {
		if ordering.Compare(current[i].listener, current[j].listener) == 0 {
			return current[i].seq < current[j].seq
		}
		return ordering.Compare(current[i].listener, current[j].listener) < 0
	})
	for _, entry := range current {
		if err := entry.listener.Handle(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// ContainerBootingEvent is published before eager singleton creation begins.
type ContainerBootingEvent struct {
	At time.Time
}

// ContainerBootedEvent is published after container boot completes.
type ContainerBootedEvent struct {
	At time.Time
}

// ContainerShuttingDownEvent is published before singleton destruction starts.
type ContainerShuttingDownEvent struct {
	At time.Time
}

// ContainerShutdownEvent is published after container shutdown completes.
type ContainerShutdownEvent struct {
	At time.Time
}

// BeanInitializedEvent is published after a singleton bean finishes initialization.
type BeanInitializedEvent struct {
	At       time.Time
	BeanName string
	BeanType reflect.Type
}

// BeanDestroyedEvent is published after a singleton bean finishes destruction.
type BeanDestroyedEvent struct {
	At       time.Time
	BeanName string
	BeanType reflect.Type
}
