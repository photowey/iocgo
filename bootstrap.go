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
	"fmt"
	"sync"
)

type Registrar struct {
	Module   string
	Register RegistrarFunc
}

type BootstrapSnapshot struct {
	BeanRegistrars    []Registrar
	StarterRegistrars []Registrar
}

type BootstrapRegistry struct {
	mu                sync.RWMutex
	beanRegistrars    map[string]Registrar
	starterRegistrars map[string]Registrar
	beanOrder         []string
	starterOrder      []string
}

func NewBootstrapRegistry() *BootstrapRegistry {
	return &BootstrapRegistry{
		beanRegistrars:    make(map[string]Registrar),
		starterRegistrars: make(map[string]Registrar),
	}
}

func (r *BootstrapRegistry) RegisterBeans(module string, fn RegistrarFunc) error {
	if module == "" {
		return fmt.Errorf("bootstrap bean registrar requires a module name")
	}
	if fn == nil {
		return fmt.Errorf("bootstrap bean registrar %q requires a register function", module)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.beanRegistrars[module]; exists {
		return fmt.Errorf("bootstrap bean registrar %q already exists", module)
	}
	r.beanRegistrars[module] = Registrar{Module: module, Register: fn}
	r.beanOrder = append(r.beanOrder, module)
	return nil
}

func (r *BootstrapRegistry) RegisterStarter(module string, fn RegistrarFunc) error {
	if module == "" {
		return fmt.Errorf("starter registrar requires a module name")
	}
	if fn == nil {
		return fmt.Errorf("starter registrar %q requires a register function", module)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.starterRegistrars[module]; exists {
		return fmt.Errorf("starter registrar %q already exists", module)
	}
	r.starterRegistrars[module] = Registrar{Module: module, Register: fn}
	r.starterOrder = append(r.starterOrder, module)
	return nil
}

func (r *BootstrapRegistry) Snapshot() BootstrapSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := BootstrapSnapshot{
		BeanRegistrars:    make([]Registrar, 0, len(r.beanOrder)),
		StarterRegistrars: make([]Registrar, 0, len(r.starterOrder)),
	}
	for _, module := range r.beanOrder {
		snapshot.BeanRegistrars = append(snapshot.BeanRegistrars, r.beanRegistrars[module])
	}
	for _, module := range r.starterOrder {
		snapshot.StarterRegistrars = append(snapshot.StarterRegistrars, r.starterRegistrars[module])
	}
	return snapshot
}

func (r *BootstrapRegistry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.beanRegistrars = make(map[string]Registrar)
	r.starterRegistrars = make(map[string]Registrar)
	r.beanOrder = nil
	r.starterOrder = nil
}

var defaultBootstrapRegistry = NewBootstrapRegistry()

func DefaultBootstrapRegistry() *BootstrapRegistry {
	return defaultBootstrapRegistry
}

func RegisterBeans(module string, fn RegistrarFunc) {
	if err := defaultBootstrapRegistry.RegisterBeans(module, fn); err != nil {
		panic(err)
	}
}

func RegisterStarter(module string, fn RegistrarFunc) {
	if err := defaultBootstrapRegistry.RegisterStarter(module, fn); err != nil {
		panic(err)
	}
}

func ResetBootstrap() {
	defaultBootstrapRegistry.Reset()
}
