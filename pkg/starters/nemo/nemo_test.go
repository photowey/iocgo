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

package nemo

import (
	"context"
	"errors"
	"strings"
	"testing"

	iocgo "github.com/photowey/iocgo"
	"github.com/photowey/nemo/pkg/collection"
)

func restoreStarter(t *testing.T) {
	t.Helper()
	iocgo.ResetBootstrap()
	Reset()
	iocgo.RegisterStarter(Module, Register)
	t.Cleanup(func() {
		iocgo.ResetBootstrap()
		Reset()
		iocgo.RegisterStarter(Module, Register)
	})
}

func TestStarterRegistersEnvironmentBean(t *testing.T) {
	restoreStarter(t)
	Configure(
		WithProperties(collection.MixedMap{
			"nemo": collection.MixedMap{
				"application": collection.MixedMap{
					"name": "iocgo-app",
				},
			},
		}),
	)

	app := iocgo.New()
	if err := app.Boot(context.Background()); err != nil {
		t.Fatalf("Boot() error = %v", err)
	}

	env, err := iocgo.Get[Environment](context.Background(), app)
	if err != nil {
		t.Fatalf("Get[Environment]() error = %v", err)
	}
	if value, ok := env.Get("nemo.application.name"); !ok || value != "iocgo-app" {
		t.Fatalf("expected application name from starter environment, got value=%v ok=%v", value, ok)
	}

	concrete, err := iocgo.Get[*StandardEnvironment](context.Background(), app, EnvironmentBeanName)
	if err != nil {
		t.Fatalf("Get[*StandardEnvironment]() error = %v", err)
	}
	if concrete == nil {
		t.Fatalf("expected concrete standard environment bean")
	}

	binder, err := iocgo.Get[*Binder](context.Background(), app, BinderBeanName)
	if err != nil {
		t.Fatalf("Get[*Binder]() error = %v", err)
	}
	if binder == nil {
		t.Fatalf("expected starter binder bean")
	}
}

func TestStarterSupportsBinding(t *testing.T) {
	restoreStarter(t)
	Configure(
		WithProperties(collection.MixedMap{
			"app": collection.MixedMap{
				"feature": collection.MixedMap{
					"enabled": "true",
					"port":    "8081",
				},
			},
		}),
	)

	type FeatureConfig struct {
		Enabled bool `binder:"enabled" required:"true"`
		Port    int  `binder:"port" default:"8081"`
	}

	app := iocgo.New()
	if err := app.Boot(context.Background()); err != nil {
		t.Fatalf("Boot() error = %v", err)
	}

	env, err := iocgo.Get[Environment](context.Background(), app)
	if err != nil {
		t.Fatalf("Get[Environment]() error = %v", err)
	}

	cfg := FeatureConfig{}
	if err := env.Bind("app.feature", &cfg); err != nil {
		t.Fatalf("Bind() error = %v", err)
	}
	if !cfg.Enabled || cfg.Port != 8081 {
		t.Fatalf("expected bound feature config, got %+v", cfg)
	}

	cfgFromHelper, err := Bind[FeatureConfig](env, "app.feature")
	if err != nil {
		t.Fatalf("Bind[T]() error = %v", err)
	}
	if !cfgFromHelper.Enabled || cfgFromHelper.Port != 8081 {
		t.Fatalf("expected helper-bound feature config, got %+v", cfgFromHelper)
	}

	binder, err := iocgo.Get[*Binder](context.Background(), app, BinderBeanName)
	if err != nil {
		t.Fatalf("Get[*Binder]() error = %v", err)
	}
	cfgFromBean := FeatureConfig{}
	if err := binder.Bind("app.feature", &cfgFromBean); err != nil {
		t.Fatalf("starter binder Bind() error = %v", err)
	}
	if !cfgFromBean.Enabled || cfgFromBean.Port != 8081 {
		t.Fatalf("expected starter binder to bind feature config, got %+v", cfgFromBean)
	}
}

func TestStarterBindingConventionCanCreateBoundBean(t *testing.T) {
	restoreStarter(t)
	Configure(
		WithProperties(collection.MixedMap{
			"app": collection.MixedMap{
				"feature": collection.MixedMap{
					"enabled": "true",
					"port":    "9090",
				},
			},
		}),
	)

	type FeatureConfig struct {
		Enabled bool `binder:"enabled" required:"true"`
		Port    int  `binder:"port" default:"8081"`
	}

	app := iocgo.New()
	if err := app.Register(
		iocgo.Define[*FeatureConfig]("featureConfig", func(ctx context.Context, resolver iocgo.Resolver) (*FeatureConfig, error) {
			bean := &FeatureConfig{}
			binder, err := iocgo.Get[*Binder](ctx, resolver, BinderBeanName)
			if err != nil {
				return nil, err
			}
			if err := binder.Bind("app.feature", bean); err != nil {
				return nil, err
			}
			return bean, nil
		}),
	); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if err := app.Boot(context.Background()); err != nil {
		t.Fatalf("Boot() error = %v", err)
	}

	cfg, err := iocgo.Get[*FeatureConfig](context.Background(), app, "featureConfig")
	if err != nil {
		t.Fatalf("Get[*FeatureConfig]() error = %v", err)
	}
	if !cfg.Enabled || cfg.Port != 9090 {
		t.Fatalf("expected bound feature config bean, got %+v", cfg)
	}
}

func TestStarterBindingSupportsDefaultValues(t *testing.T) {
	restoreStarter(t)
	Configure(
		WithProperties(collection.MixedMap{
			"app": collection.MixedMap{
				"feature": collection.MixedMap{
					"enabled": "true",
				},
			},
		}),
	)

	type FeatureConfig struct {
		Enabled bool `binder:"enabled" required:"true"`
		Port    int  `binder:"port" default:"7001"`
	}

	app := iocgo.New()
	if err := app.Boot(context.Background()); err != nil {
		t.Fatalf("Boot() error = %v", err)
	}

	env, err := iocgo.Get[Environment](context.Background(), app)
	if err != nil {
		t.Fatalf("Get[Environment]() error = %v", err)
	}
	cfg, err := Bind[FeatureConfig](env, "app.feature")
	if err != nil {
		t.Fatalf("Bind[T]() error = %v", err)
	}
	if !cfg.Enabled || cfg.Port != 7001 {
		t.Fatalf("expected default-bound config, got %+v", cfg)
	}
}

func TestStarterBindWrapsUnderlyingBindError(t *testing.T) {
	restoreStarter(t)
	Configure(
		WithProperties(collection.MixedMap{
			"app": collection.MixedMap{
				"feature": collection.MixedMap{},
			},
		}),
	)

	type FeatureConfig struct {
		Enabled bool `binder:"enabled" required:"true"`
	}

	app := iocgo.New()
	if err := app.Boot(context.Background()); err != nil {
		t.Fatalf("Boot() error = %v", err)
	}

	env, err := iocgo.Get[Environment](context.Background(), app)
	if err != nil {
		t.Fatalf("Get[Environment]() error = %v", err)
	}

	_, err = Bind[FeatureConfig](env, "app.feature")
	if err == nil {
		t.Fatalf("expected bind error")
	}
	var bindErr *BindError
	if !errors.As(err, &bindErr) {
		t.Fatalf("expected wrapped BindError, got %T", err)
	}
	if bindErr.Kind != MissingRequiredErrorKind {
		t.Fatalf("expected missing required error kind, got %s", bindErr.Kind)
	}
	if !IsBindErrorKind(err, MissingRequiredErrorKind) {
		t.Fatalf("expected IsBindErrorKind to detect missing required error")
	}
	cfgErr, ok := AsConfigurationBindError(err)
	if !ok {
		t.Fatalf("expected configuration bind error wrapper")
	}
	if cfgErr.Prefix != "app.feature" {
		t.Fatalf("expected prefix app.feature, got %q", cfgErr.Prefix)
	}
	diagnostic := FormatConfigurationBindDiagnostic(err)
	if !strings.Contains(diagnostic, "Bean:") && !strings.Contains(diagnostic, "Prefix: app.feature") {
		t.Fatalf("expected diagnostic to include starter context, got %q", diagnostic)
	}
}
