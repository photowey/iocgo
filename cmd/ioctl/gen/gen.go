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

package gen

import (
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

//
// Copyright 2018 The Kubernetes Authors.
// Copy from https://github.com/kubernetes-sigs/controller-tools/blob/master/cmd/controller-gen/main.go
//
var (
	// allGenerators maintains the list of all known generators, giving
	// them names for use on the command line.
	// each turns into a command line option,
	// and has options for output forms.
	allGenerators = map[string]genall.Generator{
		"register": crd.Generator{},
	}

	// allOutputRules defines the list of all known output rules, giving
	// them names for use on the command line.
	// Each output rule turns into two command line options:
	// - output:<generator>:<form> (per-generator output)
	// - output:<form> (default output)
	allOutputRules = map[string]genall.OutputRule{
		"dir":       genall.OutputToDirectory(""),
		"none":      genall.OutputToNothing,
		"stdout":    genall.OutputToStdout,
		"artifacts": genall.OutputArtifacts{},
	}

	// optionsRegistry contains all the marker definitions used to process command line options
	optionsRegistry = &markers.Registry{}
)

func init() {
	for genName, gen := range allGenerators {
		// make the generator options marker itself
		defn := markers.Must(markers.MakeDefinition(genName, markers.DescribesPackage, gen))
		if err := optionsRegistry.Register(defn); err != nil {
			panic(err)
		}
		if helpGiver, hasHelp := gen.(genall.HasHelp); hasHelp {
			if help := helpGiver.Help(); help != nil {
				optionsRegistry.AddHelp(defn, help)
			}
		}

		// make per-generation output rule markers
		for ruleName, rule := range allOutputRules {
			ruleMarker := markers.Must(markers.MakeDefinition(fmt.Sprintf("output:%s:%s", genName, ruleName), markers.DescribesPackage, rule))
			if err := optionsRegistry.Register(ruleMarker); err != nil {
				panic(err)
			}
			if helpGiver, hasHelp := rule.(genall.HasHelp); hasHelp {
				if help := helpGiver.Help(); help != nil {
					optionsRegistry.AddHelp(ruleMarker, help)
				}
			}
		}
	}

	// make "default output" output rule markers
	for ruleName, rule := range allOutputRules {
		ruleMarker := markers.Must(markers.MakeDefinition("output:"+ruleName, markers.DescribesPackage, rule))
		if err := optionsRegistry.Register(ruleMarker); err != nil {
			panic(err)
		}
		if helpGiver, hasHelp := rule.(genall.HasHelp); hasHelp {
			if help := helpGiver.Help(); help != nil {
				optionsRegistry.AddHelp(ruleMarker, help)
			}
		}
	}

	// add in the common options markers
	if err := genall.RegisterOptionsMarkers(optionsRegistry); err != nil {
		panic(err)
	}
}

// NoUsageError suppresses usage printing when it occurs
// (since cobra doesn't provide a good way to avoid printing
// out usage in only certain situations).
type NoUsageError struct{ error }

var Cmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate IOC bean register code.",
	Long:    "Generate IOC bean register code.",
	Example: `ioctl gen ...`,
	RunE: func(cmd *cobra.Command, rawOpts []string) error {
		fmt.Printf("prepare to generate IOC bean register code")
		switch len(rawOpts) {
		case 0:
			rawOpts = []string{"register", "paths=./..."}
		case 1:
			rawOpts = append(rawOpts, "register")
		}

		// otherwise, set up the runtime for actually running the generators
		rt, err := genall.FromOptions(optionsRegistry, rawOpts)
		if err != nil {
			return err
		}
		if len(rt.Generators) == 0 {
			return fmt.Errorf("no generators specified")
		}

		if hadErrs := rt.Run(); hadErrs {
			// don't obscure the actual error with a bunch of usage
			return NoUsageError{fmt.Errorf("not all generators ran successfully")}
		}

		return nil
	},
}
