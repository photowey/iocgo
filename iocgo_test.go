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
	"strings"
	"testing"
)

type greetingService interface {
	Message() string
}

type alphaGreeting struct {
	id int
}

func (g *alphaGreeting) Message() string { return "alpha" }

type betaGreeting struct {
	id int
}

func (g *betaGreeting) Message() string { return "beta" }

type lifecycleBean struct {
	initialized bool
	destroyed   bool
	hookTrace   []string
}

func (b *lifecycleBean) AfterPropertiesSet() {
	b.initialized = true
	b.hookTrace = append(b.hookTrace, "after-properties-set")
}

func (b *lifecycleBean) Destroy() {
	b.destroyed = true
	b.hookTrace = append(b.hookTrace, "destroy")
}

type cycleA struct{ B *cycleB }
type cycleB struct{ A *cycleA }

func TestSingletonAndPrototypeScopes(t *testing.T) {
	ResetBootstrap()
	ctx := context.Background()
	app := New()

	singletonCalls := 0
	prototypeCalls := 0

	err := app.Register(
		Define("singletonBean", func(context.Context, Resolver) (*alphaGreeting, error) {
			singletonCalls++
			return &alphaGreeting{id: singletonCalls}, nil
		}),
		Define("prototypeBean", func(context.Context, Resolver) (*betaGreeting, error) {
			prototypeCalls++
			return &betaGreeting{id: prototypeCalls}, nil
		}, WithScope(Prototype)),
	)
	if err != nil {
		t.Fatalf("register definitions: %v", err)
	}
	if err := app.Boot(ctx); err != nil {
		t.Fatalf("boot application: %v", err)
	}

	one, err := Get[*alphaGreeting](ctx, app, "singletonBean")
	if err != nil {
		t.Fatalf("get singleton first time: %v", err)
	}
	two, err := Get[*alphaGreeting](ctx, app, "singletonBean")
	if err != nil {
		t.Fatalf("get singleton second time: %v", err)
	}
	if one != two {
		t.Fatalf("expected singleton instance reuse")
	}
	if singletonCalls != 1 {
		t.Fatalf("expected singleton factory to run once, got %d", singletonCalls)
	}

	p1, err := Get[*betaGreeting](ctx, app, "prototypeBean")
	if err != nil {
		t.Fatalf("get prototype first time: %v", err)
	}
	p2, err := Get[*betaGreeting](ctx, app, "prototypeBean")
	if err != nil {
		t.Fatalf("get prototype second time: %v", err)
	}
	if p1 == p2 || p1.id == p2.id {
		t.Fatalf("expected prototype instances to differ")
	}
	if prototypeCalls != 2 {
		t.Fatalf("expected prototype factory to run twice, got %d", prototypeCalls)
	}
}

func TestPrimaryAndNamedResolution(t *testing.T) {
	ResetBootstrap()
	ctx := context.Background()
	app := New()

	err := app.Register(
		Define("alpha", func(context.Context, Resolver) (*alphaGreeting, error) {
			return &alphaGreeting{id: 1}, nil
		}, WithExposedTypes(TypeOf[greetingService]())),
		Define("beta", func(context.Context, Resolver) (*betaGreeting, error) {
			return &betaGreeting{id: 2}, nil
		}, WithExposedTypes(TypeOf[greetingService]()), WithPrimary()),
	)
	if err != nil {
		t.Fatalf("register definitions: %v", err)
	}
	if err := app.Boot(ctx); err != nil {
		t.Fatalf("boot application: %v", err)
	}

	service, err := Get[greetingService](ctx, app)
	if err != nil {
		t.Fatalf("resolve primary greeting service: %v", err)
	}
	if service.Message() != "beta" {
		t.Fatalf("expected primary bean to resolve, got %q", service.Message())
	}

	named, err := Get[greetingService](ctx, app, "alpha")
	if err != nil {
		t.Fatalf("resolve named greeting service: %v", err)
	}
	if named.Message() != "alpha" {
		t.Fatalf("expected named bean alpha, got %q", named.Message())
	}
}

func TestLifecycleHooksAndShutdown(t *testing.T) {
	ResetBootstrap()
	ctx := context.Background()
	app := New()

	err := app.Register(
		Define("lifecycle", func(context.Context, Resolver) (*lifecycleBean, error) {
			return &lifecycleBean{}, nil
		},
			WithInitHook(func(_ context.Context, bean any) error {
				lifecycle, ok := bean.(*lifecycleBean)
				if !ok {
					return fmt.Errorf("unexpected bean type %T", bean)
				}
				lifecycle.hookTrace = append(lifecycle.hookTrace, "init-hook")
				return nil
			}),
			WithDestroyHook(func(_ context.Context, bean any) error {
				lifecycle, ok := bean.(*lifecycleBean)
				if !ok {
					return fmt.Errorf("unexpected bean type %T", bean)
				}
				lifecycle.hookTrace = append(lifecycle.hookTrace, "destroy-hook")
				return nil
			}),
		),
	)
	if err != nil {
		t.Fatalf("register definitions: %v", err)
	}
	if err := app.Boot(ctx); err != nil {
		t.Fatalf("boot application: %v", err)
	}

	bean, err := Get[*lifecycleBean](ctx, app, "lifecycle")
	if err != nil {
		t.Fatalf("resolve lifecycle bean: %v", err)
	}
	if !bean.initialized {
		t.Fatalf("expected initializing lifecycle to run")
	}
	if got, want := bean.hookTrace[0], "after-properties-set"; got != want {
		t.Fatalf("expected first init trace %q, got %q", want, got)
	}
	if got, want := bean.hookTrace[1], "init-hook"; got != want {
		t.Fatalf("expected second init trace %q, got %q", want, got)
	}

	if err := app.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown application: %v", err)
	}
	if !bean.destroyed {
		t.Fatalf("expected destroy lifecycle to run")
	}
	if got, want := bean.hookTrace[len(bean.hookTrace)-1], "destroy-hook"; got != want {
		t.Fatalf("expected final destroy trace %q, got %q", want, got)
	}
}

func TestBootstrapRegistrars(t *testing.T) {
	ResetBootstrap()
	t.Cleanup(ResetBootstrap)
	ctx := context.Background()

	RegisterBeans("module-a", func(reg Registry) error {
		return reg.Register(Define("alphaBean", func(context.Context, Resolver) (*alphaGreeting, error) {
			return &alphaGreeting{id: 1}, nil
		}))
	})
	RegisterStarter("starter-a", func(reg Registry) error {
		return reg.Register(Define("betaBean", func(context.Context, Resolver) (*betaGreeting, error) {
			return &betaGreeting{id: 2}, nil
		}))
	})

	app := New()
	if err := app.Boot(ctx); err != nil {
		t.Fatalf("boot application: %v", err)
	}
	if _, err := Get[*alphaGreeting](ctx, app, "alphaBean"); err != nil {
		t.Fatalf("resolve bean registrar output: %v", err)
	}
	if _, err := Get[*betaGreeting](ctx, app, "betaBean"); err != nil {
		t.Fatalf("resolve starter registrar output: %v", err)
	}
}

func TestCycleDetection(t *testing.T) {
	ResetBootstrap()
	ctx := context.Background()
	app := New()

	defA := Define("cycleA", func(ctx context.Context, resolver Resolver) (*cycleA, error) {
		bean, err := resolver.Resolve(ctx, reflect.TypeOf((*cycleB)(nil)), "cycleB")
		if err != nil {
			return nil, err
		}
		return &cycleA{B: bean.(*cycleB)}, nil
	})
	defB := Define("cycleB", func(ctx context.Context, resolver Resolver) (*cycleB, error) {
		bean, err := resolver.Resolve(ctx, reflect.TypeOf((*cycleA)(nil)), "cycleA")
		if err != nil {
			return nil, err
		}
		return &cycleB{A: bean.(*cycleA)}, nil
	})
	if err := app.Register(defA, defB); err != nil {
		t.Fatalf("register cycle definitions: %v", err)
	}

	_, err := Get[*cycleA](ctx, app, "cycleA")
	if err == nil {
		t.Fatalf("expected cycle detection error")
	}
	if !strings.Contains(err.Error(), "bean cycle detected") {
		t.Fatalf("expected cycle error, got %q", err.Error())
	}
}
