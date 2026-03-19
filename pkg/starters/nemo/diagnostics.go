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
	"fmt"
	"strings"

	nemoapi "github.com/photowey/nemo"
)

func FormatConfigurationBindDiagnostic(err error) string {
	if err == nil {
		return "Nemo Starter Diagnostic\n  Status: no error"
	}

	cfgErr, ok := AsConfigurationBindError(err)
	if !ok {
		return fmt.Sprintf("Nemo Starter Diagnostic\n  Error: %v", err)
	}

	lines := []string{
		"Nemo Starter Diagnostic",
	}
	if cfgErr.BeanName != "" {
		lines = append(lines, fmt.Sprintf("  Bean: %s", cfgErr.BeanName))
	}
	if cfgErr.TypeName != "" {
		lines = append(lines, fmt.Sprintf("  Type: %s", cfgErr.TypeName))
	}
	if cfgErr.Prefix != "" {
		lines = append(lines, fmt.Sprintf("  Prefix: %s", cfgErr.Prefix))
	}
	if cfgErr.Cause != nil {
		lines = append(lines, fmt.Sprintf("  Cause: %v", cfgErr.Cause))
	}

	if bindErr, ok := nemoapi.AsBindError(err); ok {
		lines = append(lines, "  Underlying:")
		for _, line := range strings.Split(nemoapi.FormatBindDiagnostic(bindErr), "\n") {
			lines = append(lines, "    "+line)
		}
	}

	return strings.Join(lines, "\n")
}
