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

// +ioc:package
// +ioc:init-register
package firstpartystarterdemo

import (
	"context"

	starterNemo "github.com/photowey/iocgo/pkg/starters/nemo"
)

// +ioc:component
// +ioc:component:name=featureConfig
// +ioc:configuration-properties
// +ioc:configuration-properties:prefix=app.feature
type FeatureConfig struct {
	Enabled bool `binder:"enabled" required:"true"`
	Port    int  `binder:"port" default:"8080"`
}

// +ioc:configuration
// +ioc:configuration-properties
// +ioc:configuration-properties:prefix=app.info
type AppConfiguration struct {
	Name string `binder:"name" default:"iocgo-demo"`
}

type AppInfo struct {
	Name    string
	Feature *FeatureConfig
}

// +ioc:bean
// +ioc:bean:name=appInfo
func (c *AppConfiguration) CreateAppInfo(ctx context.Context, env starterNemo.Environment, feature *FeatureConfig) *AppInfo {
	value, _ := env.Get("nemo.application.name")
	name, _ := value.(string)
	if c.Name != "" {
		name = c.Name
	}
	return &AppInfo{
		Name:    name,
		Feature: feature,
	}
}
