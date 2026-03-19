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

// +ioc:package
// +ioc:init-register
package simple

import "context"

// +ioc:component
// +ioc:component:name=repo
// +ioc:scope=singleton
type Repo struct{}

type Service struct {
	Repo *Repo
}

// +ioc:configuration
type AppConfiguration struct{}

// +ioc:bean
// +ioc:bean:name=service
func (c *AppConfiguration) CreateService(ctx context.Context, repo *Repo) *Service {
	return &Service{Repo: repo}
}
