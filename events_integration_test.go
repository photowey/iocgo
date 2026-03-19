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
	"sync"
	"testing"

	iocgoevents "github.com/photowey/iocgo/events"
)

type userCreatedEvent struct {
	ID string
}

type listenerRecorder struct {
	mu     sync.Mutex
	events []string
}

func (r *listenerRecorder) append(value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, value)
}

func (r *listenerRecorder) snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, len(r.events))
	copy(out, r.events)
	return out
}

type orderedUserCreatedListener struct {
	order    int
	recorder *listenerRecorder
	prefix   string
}

func (l *orderedUserCreatedListener) EventType() reflect.Type {
	return reflect.TypeOf(userCreatedEvent{})
}

func (l *orderedUserCreatedListener) Handle(_ context.Context, event any) error {
	l.recorder.append(l.prefix + ":" + event.(userCreatedEvent).ID)
	return nil
}

func (l *orderedUserCreatedListener) Order() int {
	return l.order
}

func TestApplicationDiscoversListenerBeansAndPublisherBean(t *testing.T) {
	ResetBootstrap()
	ctx := context.Background()
	app := New()
	recorder := &listenerRecorder{}

	err := app.Register(
		Define[iocgoevents.Listener]("listenerA", func(context.Context, Resolver) (iocgoevents.Listener, error) {
			return &orderedUserCreatedListener{order: 10, recorder: recorder, prefix: "second"}, nil
		}),
		Define[iocgoevents.Listener]("listenerB", func(context.Context, Resolver) (iocgoevents.Listener, error) {
			return &orderedUserCreatedListener{order: 0, recorder: recorder, prefix: "first"}, nil
		}),
		Define("publisherUser", func(ctx context.Context, resolver Resolver) (*alphaGreeting, error) {
			publisher, err := Get[iocgoevents.Publisher](ctx, resolver, iocgoevents.PublisherBeanName)
			if err != nil {
				return nil, err
			}
			if err := publisher.Publish(ctx, userCreatedEvent{ID: "u-1"}); err != nil {
				return nil, err
			}
			return &alphaGreeting{id: 1}, nil
		}),
	)
	if err != nil {
		t.Fatalf("register definitions: %v", err)
	}
	if err := app.Boot(ctx); err != nil {
		t.Fatalf("boot application: %v", err)
	}

	got := recorder.snapshot()
	if len(got) != 2 || got[0] != "first:u-1" || got[1] != "second:u-1" {
		t.Fatalf("listener events = %#v", got)
	}
}

type lifecycleEventRecorder struct {
	mu      sync.Mutex
	booted  bool
	beans   []string
	destroy []string
}

func (r *lifecycleEventRecorder) EventType() reflect.Type {
	return reflect.TypeOf(iocgoevents.ContainerBootedEvent{})
}

func (r *lifecycleEventRecorder) Handle(_ context.Context, event any) error {
	switch current := event.(type) {
	case iocgoevents.ContainerBootedEvent:
		r.mu.Lock()
		r.booted = true
		r.mu.Unlock()
	case iocgoevents.BeanInitializedEvent:
		r.mu.Lock()
		r.beans = append(r.beans, current.BeanName)
		r.mu.Unlock()
	case iocgoevents.BeanDestroyedEvent:
		r.mu.Lock()
		r.destroy = append(r.destroy, current.BeanName)
		r.mu.Unlock()
	}
	return nil
}

func TestContainerAndBeanLifecycleEvents(t *testing.T) {
	ResetBootstrap()
	ctx := context.Background()
	app := New()
	recorder := &lifecycleEventRecorder{}

	err := app.Register(
		Define[iocgoevents.Listener]("bootListener", func(context.Context, Resolver) (iocgoevents.Listener, error) {
			return iocgoevents.NewListener(func(_ context.Context, event iocgoevents.ContainerBootedEvent) error {
				if err := recorder.Handle(context.Background(), event); err != nil {
					return err
				}
				return nil
			}), nil
		}),
		Define[iocgoevents.Listener]("beanInitListener", func(context.Context, Resolver) (iocgoevents.Listener, error) {
			return iocgoevents.NewListener(func(_ context.Context, event iocgoevents.BeanInitializedEvent) error {
				if err := recorder.Handle(context.Background(), event); err != nil {
					return err
				}
				return nil
			}), nil
		}),
		Define[iocgoevents.Listener]("beanDestroyListener", func(context.Context, Resolver) (iocgoevents.Listener, error) {
			return iocgoevents.NewListener(func(_ context.Context, event iocgoevents.BeanDestroyedEvent) error {
				if err := recorder.Handle(context.Background(), event); err != nil {
					return err
				}
				return nil
			}), nil
		}),
		Define("trackedBean", func(context.Context, Resolver) (*alphaGreeting, error) {
			return &alphaGreeting{id: 99}, nil
		}),
	)
	if err != nil {
		t.Fatalf("register definitions: %v", err)
	}
	if err := app.Boot(ctx); err != nil {
		t.Fatalf("boot application: %v", err)
	}
	if !recorder.booted {
		t.Fatalf("expected boot event to be observed")
	}
	if err := app.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown application: %v", err)
	}
	foundInit := false
	for _, name := range recorder.beans {
		if name == "trackedBean" {
			foundInit = true
			break
		}
	}
	if !foundInit {
		t.Fatalf("expected trackedBean initialized event, got %#v", recorder.beans)
	}
	foundDestroy := false
	for _, name := range recorder.destroy {
		if name == "trackedBean" {
			foundDestroy = true
			break
		}
	}
	if !foundDestroy {
		t.Fatalf("expected trackedBean destroyed event, got %#v", recorder.destroy)
	}
}
