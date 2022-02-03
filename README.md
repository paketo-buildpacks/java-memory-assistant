# `gcr.io/paketo-buildpacks/java-memory-assistant`

The Paketo Java Memory Assistant Buildpack is a Cloud Native Buildpack that configures the SAP Java Memory Assistant (JMA) Agent for Java applications.

## Behavior
This buildpack will participate if all the following conditions are met

* `$BP_JMA_ENABLED` is set

The buildpack will do the following:

* Contributes a Java agent to a layer and configures `$JAVA_TOOL_OPTIONS` to use it

## Configuration
| Environment Variable | Description
| -------------------- | -----------
| `$BP_JMA_ENABLED` | Whether to contribute the Java Memory Assistant Agent at build time. Defaults to `false`.
| `$BPL_JMA_ENABLED` | Whether to enable the JMA agent at runtime. Defaults to `false`.
| `$BPL_JMA_ARGS` | Configuration options for the JMA Agent. This should be a comma separated list of key/value pairs. Defaults: `check_interval=5s,log_level=ERROR,max_frequency=1/1m,heap_dump_folder=/tmp,thresholds.heap=80%`. If any custom values are specified, no default args are supplied. See Agent Configuration below for supported arguments. Example: `BPL_JMA_ARGS="check_interval=3s,heap_dump_folder=/tmp,thresholds.heap=80%"`

## Agent Configuration
| Environment Variable | Description
| -------------------- | -----------
| `heap_dump_folder=<dir>` | The folder on the container's filesystem where heap dumps are created. Default value: `/tmp`
| `thresholds.<memory area>=<value>` | (Required) This configuration allows to define thresholds for every memory area of the JVM. Thresholds can be defined in absolute percentages, e.g., 75% creates a heap dump at 75% of the selected memory area. It is also possible to specify relative increases and decreases of memory usage: for example, +5%/2m will trigger a heap dump if the particular memory area has increased by 5% or more over the last two minutes. See below to check which memory areas are supported. Thresholds can also be specified in terms of absolute values, e.g., >400MB (more than 400 MB) or <=30KB (30 KB or less); supported memory size units are KB, MB and GB. Defaults to `thresholds.heap=80%`
| `check_interval=<value>` | (Required) The interval between checks. Examples: `1s` (once a second), `3m` (every three minutes), `1h` (once every hour). Default: `5s` (check every five seconds).
| `max_frequency=<value>` | Maximum amount of heap dumps that the Java Memory Assistant is allowed to create in a given amount of time. Example: `1/30s` (no more than one heap dump every thirty seconds). The time interval is checked every time one heap dump should be created (based on the specified thresholds), and compared with the timestamps of the previously created heap dumps to make sure that the maximum frequency is not exceeded. Default: `1/1m` (one heap dump per minute)
| `log_level=<value>` | The log level used by the Java Memory Assistant. Valid values: `DEBUG`, `WARN`, `INFO`, `ERROR` and `FATAL`. Defaults to `ERROR`

Full details on the arguments supported by the agent, including the threshold memory areas available for JVM providers, can be found [here](https://github.com/SAP/java-memory-assistant).

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

