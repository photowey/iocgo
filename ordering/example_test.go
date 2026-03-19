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

package ordering_test

import (
	"fmt"

	"github.com/photowey/iocgo/ordering"
)

type priorityStep struct {
	name  string
	order int
}

func (s priorityStep) Order() int     { return s.order }
func (priorityStep) PriorityOrder()   {}
func (s priorityStep) String() string { return s.name }

type orderedStep struct {
	name  string
	order int
}

func (s orderedStep) Order() int     { return s.order }
func (s orderedStep) String() string { return s.name }

func ExampleSort() {
	steps := []fmt.Stringer{
		orderedStep{name: "ordered", order: 10},
		priorityStep{name: "priority", order: 5},
		orderedStep{name: "fallback", order: 20},
	}
	ordering.Sort(steps)
	for _, step := range steps {
		fmt.Println(step.String())
	}
	// Output:
	// priority
	// ordered
	// fallback
}
