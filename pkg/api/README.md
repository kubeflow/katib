# Katib API
This is APIs for Katib.
Manager service is an API for core component.
Suggestion service is an API for Suggestion services.
EarlyStopping service is an API for EarlyStopping services.

## Documentation
Please refer to the `api.md`:
 * [v1alpha1 documentation](./v1alpha1/gen-doc/api.md)
 * [v1alpha2 documentation](./v1alpha2/gen-doc/api.md)

## Update API and generate documents
When you want to edit the API, please only edit the corresponding `api.proto` and generate other files from it:
 * [v1alpha1/api.proto](./v1alpha1/api.proto)
 * [v1alpha2/api.proto](./v1alpha2/api.proto)

Documents are also generated from `api.proto` by [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc).
Running `build.sh` can update every file from `api.proto` and generate docs:
 * [v1alpha1/build.sh](./v1alpha1/build.sh)
 * [v1alpha2/build.sh](./v1alpha2/build.sh)
