## CLI Guide

### katib-cli

#### Usage

```
./katib-cli [options] [arguments]
This is katib cli client using cobra framework

Usage:
  katib-cli [command]

Available Commands:
  create      Create a resource from a file
  get         Display one or many resources
  help        Help about any command
  pull        Pull a resource from a file or from stdin.
  push        Push a resource from a file or from stdin.
  stop        Stop a resource

Flags:
  -h, --help            help for katib-cli
  -s, --server string   katib manager API endpoint (default "localhost:6789")

Use "katib-cli [command] --help" for more information about a command.
```

#### SubCommands
- `get`
```
list of resorces comannd can display includes: study, model

Available Commands:
  model       Display Model Info
  study       Display Study Info
```

- `create`
```
Create new resouce. YAML formats resource config are accepted.

Available Commands:
  study       Create a study from a file

Usage:
  katibcli create study [flags]
Flags:
  -f, --config string   File path of study config(required)
```

- `stop`
```
Specify resource ID or Name.

Available Commands:
  study       Stop a study
Usage:
  katibcli stop study [StudyID or StudyName]
```

- `push`
```
YAML or JSON formats are accepted.

Available Commands:
  model       Push a model Info from a file or from stdin
  study       Push a study Info from a file or from option

Usage:
push study -n [StudyName] -o [OwnerName] -d [StudyDescription]
push model -f [Path to model config]
```
