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

package nemo

import (
	"errors"
	"fmt"
)

type ConfigurationBindError struct {
	BeanName string
	TypeName string
	Prefix   string
	Cause    error
}

func (e *ConfigurationBindError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.BeanName == "" {
		return fmt.Sprintf("nemo starter bind failed for type %q at prefix %q: %v", e.TypeName, e.Prefix, e.Cause)
	}
	return fmt.Sprintf("nemo starter bind failed for bean %q (type %q, prefix %q): %v", e.BeanName, e.TypeName, e.Prefix, e.Cause)
}

func (e *ConfigurationBindError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func WrapBindError(beanName, typeName, prefix string, err error) error {
	if err == nil {
		return nil
	}
	return &ConfigurationBindError{
		BeanName: beanName,
		TypeName: typeName,
		Prefix:   prefix,
		Cause:    err,
	}
}

func AsConfigurationBindError(err error) (*ConfigurationBindError, bool) {
	var bindErr *ConfigurationBindError
	if errors.As(err, &bindErr) {
		return bindErr, true
	}
	return nil, false
}
