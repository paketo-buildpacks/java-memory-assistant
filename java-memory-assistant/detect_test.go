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

package java_memory_assistant_test

import (
	"os"
	"testing"

	java_memory_assistant "github.com/paketo-buildpacks/java-memory-assistant/java-memory-assistant"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect java_memory_assistant.Detect
	)

	it.Before(func() {
		Expect(os.Setenv("BP_JMA_ENABLED", "true")).To(Succeed())
	})

	it.After(func() {
		Expect(os.Unsetenv("BP_JMA_ENABLED")).To(Succeed())
	})

	it("passes detection", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: java_memory_assistant.PlanEntryAssistant},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: java_memory_assistant.PlanEntryAssistant},
						{Name: "jvm-application"},
					},
				},
			},
		}))
	})

	it("BP_JMA_ENABLED was not set", func() {
		Expect(os.Unsetenv("BP_JMA_ENABLED")).To(Succeed())

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass:  false,
			Plans: nil,
		}))
	})
}
