# Developer Guide

## Requirements

- Docker (17.05 or later.)

## Build from source code

You can build all images from source.

```bash
./scripts/build.sh
```

## Implement new suggestion algorithm

Suggestion API is defined as GRPC service at `API/api.proto`.
You can attach new algorithm easily.

- implement suggestion API
- make k8s service named vizier-suggestion-{ algorithm-name } and expose port 6789

And to add new suggestion service, you don't need to stop components ( vizier-core, modeldb, and anything) that are already running.
