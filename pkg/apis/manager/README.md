# Katib GRPC API

Katib offers the following APIs:

- DBManager service is an API for DB services.
- Suggestion service is an API for Suggestion services.
- EarlyStopping service is an API for EarlyStopping services.

## GRPC API documentation

See the [Katib v1beta1 API reference docs](./v1beta1/gen-doc/api.md).

## Update API and generate documents

When you want to edit the API, please only edit the corresponding `api.proto` and generate other files from it:

- [v1beta1/api.proto](./v1beta1/api.proto) for v1beta1.

We use [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) to
generate the API docs from `api.proto`.

Run `build.sh` to update every file from `api.proto` and generate the docs:

- [v1beta1/build.sh](./v1beta1/build.sh) for v1beta1.
