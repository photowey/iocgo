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

// Package ordering defines the shared ordering model used by iocgo and downstream
// frameworks.
//
// The ordering contract mirrors the Spring-style idea of:
//
//   - PriorityOrdered values execute first
//   - Ordered values execute after that
//   - unordered values preserve stable original order
//
// This package is intentionally small and reusable so event dispatch, container
// extension points, and downstream framework hooks can all share the same
// deterministic ordering semantics.
package ordering
