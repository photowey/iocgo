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
	"os"

	"github.com/photowey/iocgo/internal/constant"
	"github.com/photowey/iocgo/internal/environment/parser"
)

var _ StandardEnvironment = (*environment)(nil)

// StandardEnvironment standard implementation of Environment
type StandardEnvironment interface {
	Environment
	// ---------------------------------------------------------------- Get

	GetString(key string) (string, error)
	GetInt64(key string) (int64, error)
	GetInt32(key string) (int32, error)
	GetInt16(key string) (int16, error)
	GetInt8(key string) (int8, error)
	GetUInt64(key string) (uint64, error)
	GetUInt32(key string) (uint32, error)
	GetUInt16(key string) (uint16, error)
	GetUInt8(key string) (uint8, error)
	GetFloat64(key string) (float64, error)
	GetFloat32(key string) (float32, error)

	// ---------------------------------------------------------------- Set

	SetString(key string, value string)
	SetInt64(key string, value int64)
	SetInt32(key string, value int32)
	SetInt16(key string, value int16)
	SetInt8(key string, value int8)
	SetUInt64(key string, value uint64)
	SetUInt32(key string, value uint32)
	SetUInt16(key string, value uint16)
	SetUInt8(key string, value uint8)
	SetFloat64(key string, value float64)
	SetFloat32(key string, value float32)
}

type environment struct {
	keyDelimiter      string
	configName        string
	configType        string
	envPrefix         string
	mergeDepth        uint8
	configFiles       []string
	absPaths          []string
	searchPaths       []string
	activeProfiles    []string
	configPermissions os.FileMode
	configs           map[string]any
	defaults          map[string]any
	aliases           map[string]any
	parser            parser.StandardParser
}

func NewEnvironment(opts ...Option) StandardEnvironment {
	options := initOptions(opts...)
	env := newEnvironment(options)
	env.Init()

	return env
}

func newEnvironment(opts *Options) StandardEnvironment {
	return &environment{
		keyDelimiter:      opts.KeyDelimiter,              // .
		configName:        opts.ConfigName,                // config
		configType:        opts.ConfigType,                // yml
		envPrefix:         opts.EnvPrefix,                 // IOC_GO
		mergeDepth:        opts.MergeDepth,                // 8
		configFiles:       constant.DefaultConfigFiles,    // configs/config.yml
		absPaths:          opts.AbsPaths,                  // . config configs
		searchPaths:       opts.SearchPaths,               // . config configs
		activeProfiles:    constant.DefaultActiveProfiles, // dev
		configPermissions: os.FileMode(0x644),             // 0o644
		configs:           make(map[string]any, 0),
		defaults:          make(map[string]any, 0),
		aliases:           make(map[string]any, 0),
		parser:            parser.NewParser(),
	}
}

// ---------------------------------------------------------------- Init

func (env *environment) Init() {

	// determine configs by absPaths and searchPaths etc.
	// TODO init
}

// ---------------------------------------------------------------- Get

// ---------------------------------------------------------------- any

func (env *environment) GetProperty(key string, standBy any) (any, error) {
	return standBy, nil
}

// ---------------------------------------------------------------- string

func (env *environment) GetString(key string) (string, error) {
	return "", nil
}

// ---------------------------------------------------------------- int

func (env *environment) GetInt64(key string) (int64, error) {
	return 0, nil
}

func (env *environment) GetInt32(key string) (int32, error) {
	return 0, nil
}

func (env *environment) GetInt16(key string) (int16, error) {
	return 0, nil
}

func (env *environment) GetInt8(key string) (int8, error) {
	return 0, nil
}

// ---------------------------------------------------------------- uint

func (env *environment) GetUInt64(key string) (uint64, error) {
	return 0, nil
}

func (env *environment) GetUInt32(key string) (uint32, error) {
	return 0, nil
}

func (env *environment) GetUInt16(key string) (uint16, error) {
	return 0, nil
}

func (env *environment) GetUInt8(key string) (uint8, error) {
	return 0, nil
}

// ---------------------------------------------------------------- uint

func (env *environment) GetFloat64(key string) (float64, error) {
	return 0.0, nil
}

func (env *environment) GetFloat32(key string) (float32, error) {
	return 0.0, nil
}

// ---------------------------------------------------------------- Set

func (env *environment) SetProperty(key string, value any) {

}

func (env *environment) SetString(key string, value string) {

}

func (env *environment) SetInt64(key string, value int64) {

}

func (env *environment) SetInt32(key string, value int32) {

}

func (env *environment) SetInt16(key string, value int16) {

}

func (env *environment) SetInt8(key string, value int8) {

}

func (env *environment) SetUInt64(key string, value uint64) {

}

func (env *environment) SetUInt32(key string, value uint32) {

}

func (env *environment) SetUInt16(key string, value uint16) {

}

func (env *environment) SetUInt8(key string, value uint8) {

}

func (env *environment) SetFloat64(key string, value float64) {

}

func (env *environment) SetFloat32(key string, value float32) {

}
