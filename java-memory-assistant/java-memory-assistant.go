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
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/libpak/sherpa"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type javaMemoryAssistant struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func JavaMemoryAssistant(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (javaMemoryAssistant, libcnb.BOMEntry) {

	// Call libpak method to create a new 'contributor' which contributes our dependency to a 'Launch' layer
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return javaMemoryAssistant{LayerContributor: contributor}, entry
}

func (j javaMemoryAssistant) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {

		file := filepath.Join(layer.Path, filepath.Base(artifact.Name()))

		j.LayerContributor.Logger.Bodyf("Copying to: %s", file)
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		layer.LaunchEnvironment.Appendf("JAVA_TOOL_OPTIONS", " ",
			"-javaagent:%s", file)

		return layer, nil
	})
}

func (w javaMemoryAssistant) Name() string {
	return w.LayerContributor.LayerName()
}
