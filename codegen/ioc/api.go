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

package ioc

import (
	"fmt"

	diimpl "github.com/photowey/iocgo/cmd/ioctl/di"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type IGenerator = diimpl.IGenerator
type Generator = diimpl.Generator
type IocGoGenerationContext = diimpl.IocGoGenerationContext

type NoUsageError struct{ error }

var (
	allGenerators = map[string]genall.Generator{
		"register": Generator{},
	}

	allOutputRules = map[string]genall.OutputRule{
		"dir":       genall.OutputToDirectory(""),
		"none":      genall.OutputToNothing,
		"stdout":    genall.OutputToStdout,
		"artifacts": genall.OutputArtifacts{},
	}
)

func NewOptionsRegistry() *markers.Registry {
	registry := &markers.Registry{}
	for genName, gen := range allGenerators {
		defn := markers.Must(markers.MakeDefinition(genName, markers.DescribesPackage, gen))
		if err := registry.Register(defn); err != nil {
			panic(err)
		}
		if markerGen, ok := gen.(interface{ RegisterMarkers(*markers.Registry) error }); ok {
			if err := markerGen.RegisterMarkers(registry); err != nil {
				panic(err)
			}
		}
		if helpGiver, hasHelp := gen.(genall.HasHelp); hasHelp {
			if help := helpGiver.Help(); help != nil {
				registry.AddHelp(defn, help)
			}
		}
		for ruleName, rule := range allOutputRules {
			ruleMarker := markers.Must(markers.MakeDefinition(fmt.Sprintf("output:%s:%s", genName, ruleName), markers.DescribesPackage, rule))
			if err := registry.Register(ruleMarker); err != nil {
				panic(err)
			}
			if helpGiver, hasHelp := rule.(genall.HasHelp); hasHelp {
				if help := helpGiver.Help(); help != nil {
					registry.AddHelp(ruleMarker, help)
				}
			}
		}
	}

	for ruleName, rule := range allOutputRules {
		ruleMarker := markers.Must(markers.MakeDefinition("output:"+ruleName, markers.DescribesPackage, rule))
		if err := registry.Register(ruleMarker); err != nil {
			panic(err)
		}
		if helpGiver, hasHelp := rule.(genall.HasHelp); hasHelp {
			if help := helpGiver.Help(); help != nil {
				registry.AddHelp(ruleMarker, help)
			}
		}
	}

	if err := genall.RegisterOptionsMarkers(registry); err != nil {
		panic(err)
	}
	return registry
}

func NormalizeRawOptions(rawOpts []string) []string {
	switch len(rawOpts) {
	case 0:
		return []string{"register", "paths=./..."}
	case 1:
		return append(append([]string(nil), rawOpts...), "register")
	default:
		return append([]string(nil), rawOpts...)
	}
}

func Run(rawOpts []string) error {
	registry := NewOptionsRegistry()
	normalized := NormalizeRawOptions(rawOpts)

	rt, err := genall.FromOptions(registry, normalized)
	if err != nil {
		return err
	}
	if len(rt.Generators) == 0 {
		return fmt.Errorf("no generators specified")
	}
	if hadErrs := rt.Run(); hadErrs {
		return NoUsageError{fmt.Errorf("not all generators ran successfully")}
	}
	return nil
}
