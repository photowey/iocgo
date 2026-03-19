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
package bindcfg

// +ioc:component
// +ioc:component:name=featureConfig
// +ioc:configuration-properties
// +ioc:configuration-properties:prefix=app.feature
type FeatureConfig struct {
	Enabled bool `binder:"enabled"`
	Port    int  `binder:"port"`
}

// +ioc:configuration
// +ioc:configuration-properties
// +ioc:configuration-properties:prefix=app.info
type AppConfiguration struct {
	Name string `binder:"name"`
}
