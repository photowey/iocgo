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

package codec

var _ Registry = (*registry)(nil)

type Registry interface {
	Register(name string, handler Codec) error
	Get(name string) (Codec, error)
}

type registry struct {
}

func NewRegistry() Registry {
	return &registry{}
}

func (r registry) Register(name string, handler Codec) error {
	return nil
}

func (r registry) Get(name string) (Codec, error) {
	return nil, nil
}
