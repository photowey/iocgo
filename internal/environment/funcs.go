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

package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/photowey/iocgo/internal/constant"
)

func DetermineAbsPath(path string) string {
	if path == constant.EmptyString {
		path = constant.DotSeparator
	}
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	p, err := filepath.Abs(path)
	if err == nil {
		return filepath.Clean(p)
	}

	return constant.EmptyString
}

func DeterminePathSuffix(searchPath string) string {
	if searchPath == constant.EmptyString {
		searchPath = constant.DotSeparator
	}
	if strings.HasSuffix(searchPath, constant.PathSeparator) {
		return searchPath
	}

	return searchPath + constant.PathSeparator
}

func DetermineConfigFiles(opts *Options) []string {
	configNames := make([]string, len(opts.ActiveProfiles)+1)
	configName := PopulateConfigName(opts.ConfigName, "", opts.ConfigType)
	configNames[0] = configName
	for i, profile := range opts.ActiveProfiles {
		configNames[i+1] = PopulateConfigName(opts.ConfigName, profile, opts.ConfigType)
	}

	return configNames
}

func PopulateConfigName(configName, profile, configType string) string {
	if profile == constant.EmptyString {
		return fmt.Sprintf("%s.%s", configName, configType) // config.yml
	}
	return fmt.Sprintf("%s_%s.%s", configName, profile, configType) // config_dev.yml
}

func FileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return !stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
