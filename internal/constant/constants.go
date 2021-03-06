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

package constant

import (
	"os"
)

const (
	EmptyString   = ""
	DotSeparator  = "."
	PathSeparator = string(os.PathSeparator)
)

const (
	DefaultConfigKeyDelimiter = "."
	DefaultConfigName         = "config"
	DefaultConfigType         = "yml"
	DefaultEnvPrefix          = "IOC_GO"
)

const (
	Zero                    = 0
	DefaultMergeDepth uint8 = 8
)

var (
	DefaultConfigFiles       = []string{"configs/config.yml"}
	DefaultActiveProfiles    = []string{"dev"}
	DefaultConfigSearchPaths = []string{".", "config", "configs"}
)
