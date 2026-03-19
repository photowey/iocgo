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

import (
	"context"
	"testing"
)

type stageHook interface {
	Name() string
}

type unorderedStageHook struct {
	name string
}

func (h unorderedStageHook) Name() string { return h.name }

type orderedStageHook struct {
	name  string
	order int
}

func (h orderedStageHook) Name() string { return h.name }
func (h orderedStageHook) Order() int   { return h.order }

type priorityStageHook struct {
	name  string
	order int
}

func (h priorityStageHook) Name() string { return h.name }
func (h priorityStageHook) Order() int   { return h.order }
func (priorityStageHook) PriorityOrder() {}

func TestGetAllSortsByPriorityOrderedThenOrdered(t *testing.T) {
	ResetBootstrap()
	app := New()
	if err := app.Register(
		Define[stageHook]("unordered", func(context.Context, Resolver) (stageHook, error) {
			return unorderedStageHook{name: "unordered"}, nil
		}),
		Define[stageHook]("ordered", func(context.Context, Resolver) (stageHook, error) {
			return orderedStageHook{name: "ordered", order: 10}, nil
		}),
		Define[stageHook]("priority", func(context.Context, Resolver) (stageHook, error) {
			return priorityStageHook{name: "priority", order: 100}, nil
		}),
		Define[stageHook]("ordered-low", func(context.Context, Resolver) (stageHook, error) {
			return orderedStageHook{name: "ordered-low", order: 0}, nil
		}),
	); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if err := app.Boot(context.Background()); err != nil {
		t.Fatalf("Boot() error = %v", err)
	}
	values, err := GetAll[stageHook](context.Background(), app)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	got := make([]string, 0, len(values))
	for _, value := range values {
		got = append(got, value.Name())
	}
	want := []string{"priority", "ordered-low", "ordered", "unordered"}
	if len(got) != len(want) {
		t.Fatalf("got len %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got=%#v want=%#v", got, want)
		}
	}
}
