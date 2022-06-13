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

package executor

import (
	"context"
)

type Runnable func(ctx context.Context)
type Callable func(ctx context.Context)

type AwaitFunc func(ctx context.Context) (any, error)

type Future interface {
	Await(ctxs ...context.Context) (any, error)
}

type Executor interface {
	Execute(task Runnable, ctx context.Context) error
}

type GoroutineExecutor interface {
	Executor
	Submit(task Callable, ctx context.Context) (Future, error)
}
