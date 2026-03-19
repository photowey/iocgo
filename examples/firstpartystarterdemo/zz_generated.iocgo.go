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

package firstpartystarterdemo

import (
	"context"

	"github.com/photowey/iocgo"
	starterNemo "github.com/photowey/iocgo/pkg/starters/nemo"
)

func RegisterIocgoBeans(reg iocgo.Registry) error {
	return reg.Register(
		iocgo.Define[*FeatureConfig]("featureConfig", func(ctx context.Context, resolver iocgo.Resolver) (*FeatureConfig, error) {
			bean := &FeatureConfig{}
			binder, err := iocgo.Get[*starterNemo.Binder](ctx, resolver, starterNemo.BinderBeanName)
			if err != nil {
				return nil, err
			}
			if err := binder.Bind("app.feature", bean); err != nil {
				return nil, err
			}
			return bean, nil
		},
			iocgo.WithScope(iocgo.Singleton),
			iocgo.WithSource(iocgo.SourceInfo{Package: "github.com/photowey/iocgo/examples/firstpartystarterdemo", File: "generated", Symbol: "FeatureConfig"}),
		),
		iocgo.Define[*AppConfiguration]("appConfiguration", func(ctx context.Context, resolver iocgo.Resolver) (*AppConfiguration, error) {
			bean := &AppConfiguration{}
			binder, err := iocgo.Get[*starterNemo.Binder](ctx, resolver, starterNemo.BinderBeanName)
			if err != nil {
				return nil, err
			}
			if err := binder.Bind("app.info", bean); err != nil {
				return nil, err
			}
			return bean, nil
		},
			iocgo.WithScope(iocgo.Singleton),
			iocgo.WithSource(iocgo.SourceInfo{Package: "github.com/photowey/iocgo/examples/firstpartystarterdemo", File: "generated", Symbol: "AppConfiguration"}),
		),
		iocgo.Define[*AppInfo]("appInfo", func(ctx context.Context, resolver iocgo.Resolver) (*AppInfo, error) {
			cfg, err := iocgo.Get[*AppConfiguration](ctx, resolver, "appConfiguration")
			if err != nil {
				return nil, err
			}
			env, err := iocgo.Get[starterNemo.Environment](ctx, resolver)
			if err != nil {
				return nil, err
			}
			feature, err := iocgo.Get[*FeatureConfig](ctx, resolver, "featureConfig")
			if err != nil {
				return nil, err
			}
			bean := cfg.CreateAppInfo(ctx, env, feature)
			return bean, nil
		},
			iocgo.WithScope(iocgo.Singleton),
			iocgo.WithSource(iocgo.SourceInfo{Package: "github.com/photowey/iocgo/examples/firstpartystarterdemo", File: "generated", Symbol: "CreateAppInfo"}),
		),
	)
}

func init() {
	iocgo.RegisterBeans("github.com/photowey/iocgo/examples/firstpartystarterdemo", RegisterIocgoBeans)
}
