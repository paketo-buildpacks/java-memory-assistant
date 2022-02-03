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

package helper_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/java-memory-assistant/helper"
	"github.com/sclevine/spec"
)

func testProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		p helper.Properties
	)

	it("contributes base JMA configuration", func() {
		Expect(p.Execute()).To(Equal(map[string]string{
			"JAVA_TOOL_OPTIONS": fmt.Sprintf("-Djma.check_interval=5s -Djma.max_frequency=1/1m -Djma.heap_dump_folder=%s", os.TempDir()),
		}))
	})

	context("$BPL_JMA_ARGS is set", func() {
		it("contributes all arguments to JMA configuration", func() {
			Expect(os.Setenv("BPL_JMA_ARGS", "check_interval=10s,max_frequency=1/1m,heap_dump_folder=/tmp/,thresholds.heap=80%,log_level=DEBUG")).To(Succeed())
			Expect(p.Execute()).To(Equal(map[string]string{
				"JAVA_TOOL_OPTIONS": "-Djma.check_interval=10s -Djma.max_frequency=1/1m -Djma.heap_dump_folder=/tmp/ -Djma.thresholds.heap=80% -Djma.log_level=DEBUG"}))
		})
	})

	context("$JAVA_TOOL_OPTIONS", func() {
		it.Before(func() {
			Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
		})

		it("contributes configuration appended to existing $JAVA_TOOL_OPTIONS", func() {
			Expect(os.Setenv("BPL_JMA_ARGS", "check_interval=10s,thresholds.heap=80%")).To(Succeed())
			Expect(p.Execute()).To(Equal(map[string]string{
				"JAVA_TOOL_OPTIONS": "test-java-tool-options -Djma.check_interval=10s -Djma.thresholds.heap=80%",
			}))
		})
	})
}
