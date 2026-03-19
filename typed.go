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
)

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func Get[T any](ctx context.Context, resolver TypeResolver, name ...string) (T, error) {
	var zero T
	beanName := ""
	if len(name) > 0 {
		beanName = name[0]
	}
	bean, err := resolver.Resolve(ctx, TypeOf[T](), beanName)
	if err != nil {
		return zero, err
	}
	typed, ok := bean.(T)
	if !ok {
		return zero, fmt.Errorf("bean %q cannot be converted to %s", beanName, TypeOf[T]())
	}
	return typed, nil
}

func MustGet[T any](ctx context.Context, resolver TypeResolver, name ...string) T {
	bean, err := Get[T](ctx, resolver, name...)
	if err != nil {
		panic(err)
	}
	return bean
}

func GetAll[T any](ctx context.Context, resolver TypeResolver) ([]T, error) {
	beans, err := resolver.ResolveAll(ctx, TypeOf[T]())
	if err != nil {
		return nil, err
	}
	typed := make([]T, 0, len(beans))
	for _, bean := range beans {
		current, ok := bean.(T)
		if !ok {
			return nil, fmt.Errorf("bean cannot be converted to %s", TypeOf[T]())
		}
		typed = append(typed, current)
	}
	return typed, nil
}
