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

package ordering

import "sort"

const (
	// HighestPrecedence sorts before every other ordered value.
	HighestPrecedence = -1 << 31
	// LowestPrecedence sorts after every other ordered value.
	LowestPrecedence = 1<<31 - 1
)

// Ordered provides a stable numeric order contract for extension points.
// Lower values execute first.
type Ordered interface {
	Order() int
}

// PriorityOrdered marks an Ordered value as part of the highest-priority
// ordering tier. Priority-ordered values always sort ahead of values that only
// implement Ordered.
type PriorityOrdered interface {
	Ordered
	PriorityOrder()
}

// metadata captures the comparable ordering tier and sort value for a candidate.
type metadata struct {
	tier  int
	order int
}

// MetadataOf returns ordering metadata for a candidate. The bool result reports
// whether the candidate participates in the shared ordering contract.
func MetadataOf(value any) (metadata, bool) {
	switch typed := value.(type) {
	case PriorityOrdered:
		return metadata{tier: 0, order: typed.Order()}, true
	case Ordered:
		return metadata{tier: 1, order: typed.Order()}, true
	default:
		return metadata{}, false
	}
}

// Compare reports relative ordering between two values using the shared
// ordering model.
//
// It returns:
//   - a negative value when left should sort before right
//   - zero when both values belong to the same ordering position
//   - a positive value when left should sort after right
func Compare(left, right any) int {
	leftMeta, leftOrdered := MetadataOf(left)
	rightMeta, rightOrdered := MetadataOf(right)

	switch {
	case leftOrdered && rightOrdered:
		if leftMeta.tier != rightMeta.tier {
			return leftMeta.tier - rightMeta.tier
		}
		return leftMeta.order - rightMeta.order
	case leftOrdered:
		return -1
	case rightOrdered:
		return 1
	default:
		return 0
	}
}

// Sort reorders values in-place according to the shared ordering contract while
// preserving stable original order for ties and unordered values.
func Sort[T any](values []T) {
	sort.SliceStable(values, func(i, j int) bool {
		return Compare(values[i], values[j]) < 0
	})
}

// SortAny reorders a slice of dynamic values in-place using the shared ordering
// contract while preserving stable original order for ties and unordered values.
func SortAny(values []any) {
	sort.SliceStable(values, func(i, j int) bool {
		return Compare(values[i], values[j]) < 0
	})
}
