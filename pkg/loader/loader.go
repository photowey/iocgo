/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package loader

import (
	"github.com/photowey/iocgo/pkg/codec"
)

const (
	DefaultConfigName = "config"
	DefaultConfigType = "yaml"
)

var (
	DefaultConfigSearchPaths = []string{".", "config", "configs"}
)

var _ Binder = (*binder)(nil)
var _ Loader = (*loaderx)(nil)

type Binder interface {
	Bind(prefix string, dst any) error
}

type binder struct {
}

func NewBinder() Binder {
	return &binder{}
}

func (bx binder) Bind(prefix string, dst any) error {

	return nil
}

type Loader interface {
	Bind(fileName, fileType string, dst any, searchPath ...string) error
	BindStruct(prefix string, dst any) error
	Load(fileName, fileType string, searchPath ...string) error
}

type loaderx struct {
	registry codec.Registry
}

func NewLoader() Loader {
	return &loaderx{
		registry: codec.NewRegistry(),
	}
}

func (lx loaderx) Bind(fileName, fileType string, dst any, searchPath ...string) error {
	return nil
}

func (lx loaderx) BindStruct(prefix string, dst any) error {
	return nil
}

func (lx loaderx) Load(fileName, fileType string, searchPath ...string) error {
	return nil
}
