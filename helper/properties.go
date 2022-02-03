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

package helper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/libpak/effect"

	"github.com/paketo-buildpacks/libpak/sherpa"

	"github.com/paketo-buildpacks/libpak/bard"
)

type Properties struct {
	Logger   bard.Logger
	Executor effect.Executor
}

func (p Properties) Execute() (map[string]string, error) {
	if val := sherpa.ResolveBool("BPL_JMA_ENABLED"); !val {
		return nil, nil
	}

	var argList string

	if argList = sherpa.GetEnvWithDefault("BPL_JMA_ARGS", ""); argList == "" {
		argList = fmt.Sprintf("-Djma.check_interval=5s -Djma.log_level=ERROR -Djma.max_frequency=1/1m -Djma.heap_dump_folder=%s -Djma.thresholds.heap=80%%", filepath.Join(os.TempDir(), "jma"))
	} else {

		// To allow simple args in BPL_JMA_ARGS, append the "-Djma." required for all configs
		var runtimeArgs []string
		for _, arg := range strings.Split(argList, ",") {
			runtimeArgs = append(runtimeArgs, fmt.Sprintf("-Djma.%s", arg))
		}
		argList = strings.Join(runtimeArgs, " ")
	}

	// Java 9+ requires an extra JVM arg to allow the creation of Heap Dumps
	if jv, err := p.javaVersion(); err == nil {
		if jv != "8" { //
			argList = fmt.Sprintf("--add-opens=jdk.management/com.sun.management.internal=ALL-UNNAMED %s", argList)
		}
	} else {
		return nil, fmt.Errorf("error checking Java version: %w", err)
	}

	p.Logger.Infof("Enabling Java Memory Assistant with args: %s", argList)

	opts := sherpa.AppendToEnvVar("JAVA_TOOL_OPTIONS", " ", argList)

	return map[string]string{"JAVA_TOOL_OPTIONS": opts}, nil
}

func (p Properties) javaVersion() (string, error) {
	buf := &bytes.Buffer{}

	path := sherpa.GetEnvWithDefault("PATH", "")

	if err := p.Executor.Execute(effect.Execution{
		Command: "java",
		Args:    []string{"-version"},
		Stdout:  buf,
		Stderr:  buf,
		Env:     []string{strings.Join([]string{"PATH", path}, ":")},
	}); err != nil {
		return "", fmt.Errorf("unable to check Java version, error: %s \n%w", strings.TrimSpace(buf.String()), err)
	}

	output := strings.Split(strings.TrimSpace(buf.String()), " ")

	//strip quotes e.g. java version "1.8.0_281"
	ver := strings.Split(strings.Trim(output[2], "\""), ".")[0]
	if ver == "1" {
		ver = "8"
	}

	return ver, nil
}
