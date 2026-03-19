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
	"fmt"
	"sync"

	iocgo "github.com/photowey/iocgo"
	nemoapi "github.com/photowey/nemo"
)

const (
	Module              = "github.com/photowey/iocgo/pkg/starters/nemo"
	EnvironmentBeanName = "nemoEnvironment"
	BinderBeanName      = "nemoBinder"
)

type Environment = nemoapi.Environment
type StandardEnvironment = nemoapi.StandardEnvironment
type Option = nemoapi.Option
type PropertySource = nemoapi.PropertySource
type PropertySources = nemoapi.PropertySources
type SuccessThreshold = nemoapi.SuccessThreshold
type MixedMap = nemoapi.MixedMap
type BindError = nemoapi.BindError
type ErrorKind = nemoapi.ErrorKind

type Binder struct {
	env Environment
}

func (b *Binder) Bind(prefix string, target any) error {
	return WrapBindError("", typeNameOfTarget(target), prefix, b.env.Bind(prefix, target))
}

var (
	NoneSuccessThreshold      = nemoapi.NoneSuccessThreshold
	AnyoneSuccessThreshold    = nemoapi.AnyoneSuccessThreshold
	AllSuccessThreshold       = nemoapi.AllSuccessThreshold
	InvalidTargetErrorKind    = nemoapi.InvalidTargetErrorKind
	InvalidTagErrorKind       = nemoapi.InvalidTagErrorKind
	UnsettableFieldErrorKind  = nemoapi.UnsettableFieldErrorKind
	MissingRequiredErrorKind  = nemoapi.MissingRequiredErrorKind
	UnsupportedTypeErrorKind  = nemoapi.UnsupportedTypeErrorKind
	ConversionFailedErrorKind = nemoapi.ConversionFailedErrorKind
)

var (
	optionsMu sync.RWMutex
	options   = defaultOptions()
)

func init() {
	iocgo.RegisterStarter(Module, Register)
}

func defaultOptions() []Option {
	return []Option{
		nemoapi.WithSearchPaths(".", "resources", "config", "configs"),
	}
}

func DefaultOptions() []Option {
	return cloneOptions(defaultOptions())
}

func Configure(opts ...Option) {
	optionsMu.Lock()
	defer optionsMu.Unlock()
	if len(opts) == 0 {
		options = defaultOptions()
		return
	}
	options = cloneOptions(opts)
}

func Reset() {
	Configure()
}

func Register(reg iocgo.Registry) error {
	return reg.Register(
		iocgo.Define[*nemoapi.StandardEnvironment](EnvironmentBeanName, func(context.Context, iocgo.Resolver) (*nemoapi.StandardEnvironment, error) {
			env := nemoapi.New().(*nemoapi.StandardEnvironment)
			if err := env.Start(currentOptions()...); err != nil {
				return nil, err
			}
			return env, nil
		},
			iocgo.WithExposedTypes(iocgo.TypeOf[nemoapi.Environment]()),
			iocgo.WithDestroyHook(func(_ context.Context, bean any) error {
				return bean.(*nemoapi.StandardEnvironment).Destroy()
			}),
		),
		iocgo.Define[*Binder](BinderBeanName, func(ctx context.Context, resolver iocgo.Resolver) (*Binder, error) {
			env, err := iocgo.Get[Environment](ctx, resolver, EnvironmentBeanName)
			if err != nil {
				return nil, err
			}
			return &Binder{env: env}, nil
		}),
	)
}

func currentOptions() []Option {
	optionsMu.RLock()
	defer optionsMu.RUnlock()
	return cloneOptions(options)
}

func cloneOptions(opts []Option) []Option {
	if len(opts) == 0 {
		return nil
	}
	cloned := make([]Option, 0, len(opts))
	cloned = append(cloned, opts...)
	return cloned
}

func WithAbsolutePaths(absolutePaths ...string) Option {
	return nemoapi.WithAbsolutePaths(absolutePaths...)
}

func WithConfigNames(configNames ...string) Option {
	return nemoapi.WithConfigNames(configNames...)
}

func WithConfigTypes(configTypes ...string) Option {
	return nemoapi.WithConfigTypes(configTypes...)
}

func WithSearchPaths(searchPaths ...string) Option {
	return nemoapi.WithSearchPaths(searchPaths...)
}

func WithProfiles(profiles ...string) Option {
	return nemoapi.WithProfiles(profiles...)
}

func WithSources(sources ...PropertySource) Option {
	return nemoapi.WithSources(sources...)
}

func WithProperties(properties MixedMap) Option {
	return nemoapi.WithProperties(properties)
}

func WithThreshold(threshold SuccessThreshold) Option {
	return nemoapi.WithThreshold(threshold)
}

func Bind[T any](env Environment, prefix string) (T, error) {
	var target T
	err := env.Bind(prefix, &target)
	err = WrapBindError("", typeNameOfTarget(target), prefix, err)
	return target, err
}

func MustBind[T any](env Environment, prefix string) T {
	target, err := Bind[T](env, prefix)
	if err != nil {
		panic(err)
	}
	return target
}

func AsBindError(err error) (*BindError, bool) {
	var bindErr *BindError
	if errors.As(err, &bindErr) {
		return bindErr, true
	}
	return nil, false
}

func IsBindErrorKind(err error, kind ErrorKind) bool {
	bindErr, ok := AsBindError(err)
	return ok && bindErr.Kind == kind
}

func typeNameOfTarget(target any) string {
	return fmt.Sprintf("%T", target)
}
