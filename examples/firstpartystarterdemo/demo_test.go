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

package firstpartystarterdemo

import (
	"context"
	"testing"

	iocgo "github.com/photowey/iocgo"
	starterNemo "github.com/photowey/iocgo/pkg/starters/nemo"
)

func TestFirstPartyStarterDemo(t *testing.T) {
	iocgo.ResetBootstrap()
	starterNemo.Reset()
	iocgo.RegisterStarter(starterNemo.Module, starterNemo.Register)
	iocgo.RegisterBeans("github.com/photowey/iocgo/examples/firstpartystarterdemo", RegisterIocgoBeans)
	t.Cleanup(func() {
		iocgo.ResetBootstrap()
		starterNemo.Reset()
		iocgo.RegisterStarter(starterNemo.Module, starterNemo.Register)
	})

	starterNemo.Configure(
		starterNemo.WithSearchPaths("configs"),
		starterNemo.WithProfiles("demo"),
	)

	app := iocgo.New()
	if err := app.Boot(context.Background()); err != nil {
		t.Fatalf("Boot() error = %v", err)
	}

	info, err := iocgo.Get[*AppInfo](context.Background(), app, "appInfo")
	if err != nil {
		t.Fatalf("Get[*AppInfo]() error = %v", err)
	}
	if info.Name != "starter-demo-config" {
		t.Fatalf("expected config-bound name override, got %q", info.Name)
	}
	if info.Feature == nil || !info.Feature.Enabled || info.Feature.Port != 8080 {
		t.Fatalf("expected bound feature config with default port, got %+v", info.Feature)
	}
}
