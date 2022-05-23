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

package beandefinition_test

import (
	"github.com/photowey/iocgo/internal/scope"
)

//
// the package of IOC bean definition.
//

// +ioc:autowired:beandefinition  // mark a beandefinition

type DumpyBeanDefinition struct {
	Name              string      // bean name
	Scope             scope.Scope // scope of bean
	AutowireCandidate bool        // autowire enabled
	Configuration     string      // configuration struct
	FactoryFunc       string      // factory func
	InitFunc          string      // init func
	DestroyFunc       string      // destroy func
	Description       string      // description of bean
}
