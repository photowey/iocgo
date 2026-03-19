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

package events

import (
	"context"
	"testing"
)

type userCreated struct {
	ID string
}

func TestSyncDispatcherPublishesInOrder(t *testing.T) {
	order := make([]string, 0, 2)
	dispatcher := NewSyncDispatcher(
		NewListener(func(_ context.Context, event userCreated) error {
			order = append(order, "second:"+event.ID)
			return nil
		}, 10),
		NewListener(func(_ context.Context, event userCreated) error {
			order = append(order, "first:"+event.ID)
			return nil
		}, 0),
	)
	if err := dispatcher.Publish(context.Background(), userCreated{ID: "u-1"}); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if len(order) != 2 || order[0] != "first:u-1" || order[1] != "second:u-1" {
		t.Fatalf("order = %#v", order)
	}
}
