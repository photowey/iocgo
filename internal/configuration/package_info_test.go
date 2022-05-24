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

package configuration_test

import (
	"context"
)

//
// configuration package
//

// +ioc:autowired:factorybean // mark a bean is a factory bean

type Person struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

// +ioc:autowired:configuration // mark a configuration struct -> @Configuration
type PersonConfiguration struct {
}

// +ioc:autowired:factoryfunc // mark a func is configuration factory func -> @Bean
// +ioc:autowired:scope=singleton // mark a bean scope is singleton
// +ioc:autowired:component=personSingleton // The id of the bean defined in the ioc

func (config *PersonConfiguration) CreatePerson(ctx context.Context) Person {
	return Person{
		ID:   9527,
		Name: "photowey",
		Age:  18,
	}
}

// +ioc:autowired:factoryfunc
// +ioc:autowired:scope=singleton // The id of the bean default is factory-func name,such as: PersonSingleton

func (config *PersonConfiguration) PersonSingleton(ctx context.Context) Person {
	return Person{
		ID:   9527,
		Name: "photowey",
		Age:  18,
	}
}
