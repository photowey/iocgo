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

// Package iocgo provides an annotation-first, static IoC container for Go.
//
// The package is intentionally centered on a small runtime API because the main
// developer experience comes from:
//
//   - generated bean registration
//   - deterministic container boot
//   - predictable lifecycle and dependency resolution
//
// This keeps framework behavior explicit and testable instead of hiding it
// behind runtime package scanning or implicit reflection-heavy discovery.
package iocgo
