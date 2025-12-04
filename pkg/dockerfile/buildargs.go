/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dockerfile

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

type BuildArgs struct {
	// args are the build args passed via CLI
	args map[string]*string
	// allowed are the args specified in the Dockerfile via ARG instruction
	allowed map[string]*string
}

func NewBuildArgs(args []string) *BuildArgs {
	argsFromOptions := make(map[string]*string)
	for _, a := range args {
		s := strings.SplitN(a, "=", 2)
		if len(s) == 1 {
			argsFromOptions[s[0]] = nil
		} else {
			argsFromOptions[s[0]] = &s[1]
		}
	}
	return &BuildArgs{
		args:    argsFromOptions,
		allowed: make(map[string]*string),
	}
}

func (b *BuildArgs) Clone() *BuildArgs {
	cloneArgs := make(map[string]*string)
	for k, v := range b.args {
		cloneArgs[k] = v
	}
	cloneAllowed := make(map[string]*string)
	for k, v := range b.allowed {
		cloneAllowed[k] = v
	}
	return &BuildArgs{
		args:    cloneArgs,
		allowed: cloneAllowed,
	}
}

// ReplacementEnvs returns a list of filtered environment variables
func (b *BuildArgs) ReplacementEnvs(envs []string) []string {
	// Ensure that we operate on a new array and do not modify the underlying array
	resultEnv := make([]string, len(envs))
	copy(resultEnv, envs)
	filtered := b.FilterAllowed(envs)
	// Disable makezero linter, since the previous make is paired with a same sized copy
	return append(resultEnv, filtered...) //nolint:makezero
}

// FilterAllowed returns the list of allowed build args that should be added to the environment
func (b *BuildArgs) FilterAllowed(envs []string) []string {
	var filtered []string
	for key, defaultVal := range b.allowed {
		// 1. Check if passed via CLI
		if val, ok := b.args[key]; ok && val != nil {
			filtered = append(filtered, key+"="+*val)
			continue
		}
		// 2. Check if it has a default value from ARG instruction
		if defaultVal != nil {
			filtered = append(filtered, key+"="+*defaultVal)
			continue
		}
		// 3. Check if it exists in the current environment (inherited)
		// This is for ARG foo (without default), it picks up foo from env if set
		for _, env := range envs {
			s := strings.SplitN(env, "=", 2)
			if len(s) == 2 && s[0] == key {
				filtered = append(filtered, env)
				break
			}
		}
	}
	return filtered
}

// AddMetaArgs adds the supplied args map to b's allowedMetaArgs
func (b *BuildArgs) AddMetaArgs(metaArgs []instructions.ArgCommand) {
	for _, marg := range metaArgs {
		for _, arg := range marg.Args {
			v := arg.Value
			b.AddMetaArg(arg.Key, v)
		}
	}
}

func (b *BuildArgs) AddMetaArg(key string, value *string) {
	b.allowed[key] = value
}

// AddArg adds a new build arg or updates an existing one
func (b *BuildArgs) AddArg(key string, value *string) {
	b.allowed[key] = value
}

// GetAllAllowed returns the map of build args passed via CLI
func (b *BuildArgs) GetAllAllowed() map[string]string {
	m := make(map[string]string)
	for k, v := range b.args {
		if v != nil {
			m[k] = *v
		}
	}
	return m
}

// GetAllMeta returns the map of args defined in Dockerfile
func (b *BuildArgs) GetAllMeta() map[string]string {
	m := make(map[string]string)
	for k, v := range b.allowed {
		if v != nil {
			m[k] = *v
		}
	}
	return m
}
