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

package version

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	Commit    = "unknown"
	BuildTime = "unknown"
)

var Cmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), Summary())
		if err != nil {
			return
		}
	},
}

// Now returns the normalized CLI version string.
func Now() string {
	value := strings.TrimSpace(Version)
	switch {
	case value == "", value == "unknown":
		return "(unknown)"
	case value == "dev":
		return value
	case strings.HasPrefix(value, "v"):
		return value
	default:
		return "v" + value
	}
}

// Summary returns the CLI version metadata summary.
func Summary() string {
	version := Now()
	commit := normalizeValue(Commit, "unknown")
	buildTime := normalizeValue(BuildTime, "unknown")
	if commit == "unknown" && buildTime == "unknown" {
		return version
	}

	return fmt.Sprintf("%s (commit %s, built %s)", version, commit, buildTime)
}

// MainVersion returns the version of the main module
func MainVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil || info.Main.Version == "" {
		// binary has not been built with module support or doesn't contain a version.
		return "(unknown)"
	}
	return info.Main.Version
}

// Print prints the CLI version summary on stdout.
func Print() {
	fmt.Printf("Version: %s\n", Summary())
}

func normalizeValue(value, fallback string) string {
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return fallback
	}

	return normalized
}
