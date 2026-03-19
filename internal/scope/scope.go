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

package scope

import "fmt"

type Scope uint8

const (
	Singleton Scope = iota + 1
	Prototype
)

func (s Scope) String() string {
	switch s {
	case Singleton:
		return "singleton"
	case Prototype:
		return "prototype"
	default:
		return "unknown"
	}
}

func (s Scope) Valid() bool {
	return s == Singleton || s == Prototype
}

func Parse(src string) (Scope, error) {
	switch src {
	case "", "singleton":
		return Singleton, nil
	case "prototype":
		return Prototype, nil
	default:
		return 0, fmt.Errorf("unsupported scope %q", src)
	}
}
