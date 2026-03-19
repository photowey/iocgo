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

package gen

import (
	"fmt"

	publicioc "github.com/photowey/iocgo/codegen/ioc"
	"github.com/spf13/cobra"
)

type NoUsageError = publicioc.NoUsageError

var Cmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate IOC bean register code.",
	Long:    "Generate IOC bean register code.",
	Example: `ioctl gen ...`,
	RunE: func(cmd *cobra.Command, rawOpts []string) error {
		fmt.Printf("prepare to generate IOC bean register code")
		if err := publicioc.Run(rawOpts); err != nil {
			if _, ok := err.(publicioc.NoUsageError); ok {
				return err
			}
			return err
		}
		return nil
	},
}
