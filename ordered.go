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

import "github.com/photowey/iocgo/ordering"

// Ordered is the shared extension ordering contract.
type Ordered = ordering.Ordered

// PriorityOrdered marks a value as part of the highest-priority ordering tier.
type PriorityOrdered = ordering.PriorityOrdered

const (
	// HighestPrecedence sorts before every other ordered value.
	HighestPrecedence = ordering.HighestPrecedence
	// LowestPrecedence sorts after every other ordered value.
	LowestPrecedence = ordering.LowestPrecedence
)

// SortByOrder reorders values in-place using the shared PriorityOrdered /
// Ordered model while preserving stable original order for ties.
func SortByOrder[T any](values []T) {
	ordering.Sort(values)
}
