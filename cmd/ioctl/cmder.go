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

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/photowey/iocgo/cmd/ioctl/gen"
	"github.com/photowey/iocgo/cmd/ioctl/version"
	"github.com/spf13/cobra"
)

var (
	App     Cmder
	rootCmd = &cobra.Command{
		Use: "ioctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to use ioctl...")
		},
	}
)

func init() {
	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(gen.Cmd)
}

type Cmder struct{}

func (app Cmder) Run() {
	if err := rootCmd.Execute(); err != nil {
		if _, noUsage := err.(gen.NoUsageError); !noUsage {
			// print the usage unless we suppressed it
			if err := rootCmd.Usage(); err != nil {
				panic(err)
			}
		}
		fmt.Fprintf(rootCmd.OutOrStderr(),
			"run `%[1]s %[2]s -w` to see all available markers, or `%[1]s %[2]s -h` for usage\n",
			rootCmd.CalledAs(), strings.Join(os.Args[1:], " "))

		os.Exit(1)
	}
}
