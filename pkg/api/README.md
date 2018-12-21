# Katib API
This is APIs for Katib.
Manager service is an API for core component.
Suggestion  service is an API for Suggestion services.
EarlyStopping service is an API for EarlyStopping services.

## Documentation
Please refer to [api.md](./gen-doc/api.md).

## Update API and generatie documents
When you want to edit API, please only edit [api.proto](./api.proto) and generate other files from it.
Documents are also generated from [api.proto](./api.proto) by [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc).
By [build.sh](./build.sh), you can update every files from api.proto and generate docs.
