# Katib GRPC API

Katib offers the following APIs:

- DBManager service is an API for DB services.
- Suggestion service is an API for Suggestion services.
- EarlyStopping service is an API for EarlyStopping services.

## GRPC API documentation

See the [Katib v1beta1 API reference docs](https://www.kubeflow.org/docs/reference/katib/v1beta1/katib/).

## Update API and generate documents

When you want to edit the API, please only edit the corresponding `api.proto` and generate other files from it:

- [v1beta1/api.proto](./v1beta1/api.proto) for v1beta1.

We use [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) to
generate the API docs from `api.proto`.

Run `build.sh` to update every file from `api.proto` and generate the docs:

- [v1beta1/build.sh](./v1beta1/build.sh) for v1beta1.

After running `build.sh`, follow these steps to update the docs:

1. Copy the updated content from your generated file
   `pkg/apis/manager/<version>/gen-doc/api.md` to the doc page in the
   `kubeflow/website` repository:
   `kubeflow/website/blob/master/content/docs/reference/katib/<version>/katib.md`.
1. Create a PR in the `kubeflow/website` repository.
   (See [example PR](https://github.com/kubeflow/website/pull/1531).)
