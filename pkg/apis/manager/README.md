# Katib API
Katib offers the following APIs:
* Manager service is an API for core component.
* Suggestion service is an API for Suggestion services.
* EarlyStopping service is an API for EarlyStopping services.

## Documentation
See the [Katib API reference docs](https://www.kubeflow.org/docs/reference/katib/).

## Update API and generate documents
When you want to edit the API, please only edit the corresponding `api.proto` and generate other files from it:
 * [v1alpha3/api.proto](./v1alpha3/api.proto)

We use [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) to
generate the API docs from `api.proto`.

Run `build.sh` to update every file from `api.proto` and generate the docs:
 * [v1alpha3/build.sh](./v1alpha3/build.sh)

After running `build.sh`, follow these steps to update the docs: 

1. Copy the updated content from your generated file
  `pkg/apis/manager/v1alpha3/gen-doc/api.md` to the doc page in the 
  `kubeflow/website` repository:
  `kubeflow/website/blob/master/content/docs/reference/katib/v1alpha3/katib.md`.
1. Create a PR in the `kubeflow/website` repository. 
  (See [example PR](https://github.com/kubeflow/website/pull/1531).)
