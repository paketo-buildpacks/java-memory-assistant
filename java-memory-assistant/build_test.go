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
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	jma "github.com/paketo-buildpacks/java-memory-assistant/java-memory-assistant"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {

	var (
		Expect = NewWithT(t).Expect
		build  jma.Build
		ctx    libcnb.BuildContext
	)

	it("contributes Java Memory Assistant agent for API <= 0.6", func() {
		ctx.Buildpack.API = "0.6"
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "java-memory-assistant"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "java-memory-assistant",
					"version": "1.0.0",
					"stacks":  []interface{}{"io.buildpacks.stacks.bionic"},
				},
			},
		}
		ctx.StackID = "io.buildpacks.stacks.bionic"

		result, err := build.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("java-memory-assistant"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.BOM.Entries).To(HaveLen(2))
	})

	it("contributes Java Memory Assistant agent", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "java-memory-assistant"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "java-memory-assistant",
					"version": "1.0.0",
					"stacks":  []interface{}{"io.buildpacks.stacks.bionic"},
					"cpes":    []string{"cpe:2.3:a:java-memory-assistant:java-memory-assistant:0.5.0:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/java-memory-assistant@1.17.1?arch=amd64",
				},
			},
		}
		ctx.StackID = "io.buildpacks.stacks.bionic"

		result, err := build.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("java-memory-assistant"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
		Expect(result.BOM.Entries).To(HaveLen(2))
	})

}
