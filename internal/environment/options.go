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
	"github.com/photowey/iocgo/internal/constant"
	"github.com/photowey/iocgo/internal/stringz"
)

type Option func(opts *Options)

type Options struct {
	KeyDelimiter   string
	ConfigName     string
	ConfigType     string
	EnvPrefix      string
	MergeDepth     uint8
	AbsPaths       []string
	SearchPaths    []string
	ActiveProfiles []string
}

func (opts *Options) validate() {
	if stringz.IsBlankString(opts.KeyDelimiter) {
		opts.KeyDelimiter = constant.DefaultConfigKeyDelimiter
	}
	if stringz.IsBlankString(opts.ConfigName) {
		opts.ConfigName = constant.DefaultConfigName
	}
	if stringz.IsBlankString(opts.ConfigType) {
		opts.ConfigType = constant.DefaultConfigType
	}
	if stringz.IsBlankString(opts.EnvPrefix) {
		opts.EnvPrefix = constant.DefaultEnvPrefix
	}
	if opts.MergeDepth == constant.Zero {
		opts.MergeDepth = constant.DefaultMergeDepth
	}
	if stringz.IsEmptyStringSlice(opts.SearchPaths) {
		opts.SearchPaths = constant.DefaultConfigSearchPaths
	}
	if stringz.IsEmptyStringSlice(opts.ActiveProfiles) {
		opts.ActiveProfiles = constant.DefaultActiveProfiles
	}
}

func NewOptions() *Options {
	return &Options{
		AbsPaths:       make([]string, 0),
		SearchPaths:    make([]string, 0),
		ActiveProfiles: make([]string, 0),
	}
}

func initOptions(opts ...Option) *Options {
	options := NewOptions()
	for _, opt := range opts {
		opt(options)
	}
	options.validate()

	return options
}

// ----------------------------------------------------------------

func WithKeyDelimiter(keyDelimiter string) Option {
	return func(opts *Options) {
		opts.KeyDelimiter = keyDelimiter
	}
}

func WithConfigName(configName string) Option {
	return func(opts *Options) {
		opts.ConfigName = configName
	}
}

func WithConfigType(configType string) Option {
	return func(opts *Options) {
		opts.ConfigType = configType
	}
}

func WithEnvPrefix(envPrefix string) Option {
	return func(opts *Options) {
		opts.EnvPrefix = envPrefix
	}
}

func WithMergeDepth(mergeDepth uint8) Option {
	return func(opts *Options) {
		opts.MergeDepth = mergeDepth
	}
}

func WithAbsPaths(absPaths ...string) Option {
	return func(opts *Options) {
		opts.AbsPaths = absPaths
	}
}

func WithSearchPaths(searchPaths ...string) Option {
	return func(opts *Options) {
		opts.SearchPaths = searchPaths
	}
}

func WithActiveProfiles(activeProfiles ...string) Option {
	return func(opts *Options) {
		opts.ActiveProfiles = activeProfiles
	}
}

// ----------------------------------------------------------------
