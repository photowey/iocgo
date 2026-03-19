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

package iocgo_test

import (
	"context"
	"fmt"

	iocgo "github.com/photowey/iocgo"
)

type ExampleService struct{}

func (ExampleService) Message() string { return "hello" }

func ExampleNew() {
	app := iocgo.New()
	_ = app.Register(
		iocgo.Define[*ExampleService]("exampleService", func(context.Context, iocgo.Resolver) (*ExampleService, error) {
			return &ExampleService{}, nil
		}),
	)
	_ = app.Boot(context.Background())
	service, _ := iocgo.Get[*ExampleService](context.Background(), app, "exampleService")
	fmt.Println(service.Message())
	// Output:
	// hello
}
