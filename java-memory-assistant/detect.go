/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java_memory_assistant

import (
	"github.com/paketo-buildpacks/libpak/sherpa"

	"github.com/buildpacks/libcnb"
)

const (
	PlanEntryAssistant = "java-memory-assistant"
)

type Detect struct {
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {

	if val := sherpa.ResolveBool("BP_JMA_ENABLED"); !val {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				// Indicates that our Buildpack 'provides' a dependency called 'java-memory-assistant'
				Provides: []libcnb.BuildPlanProvide{
					{Name: PlanEntryAssistant},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: PlanEntryAssistant},
					{Name: "jvm-application"},
				},
			},
		},
	}, nil
}
