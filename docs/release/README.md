# Release the Katib Project

This is the instruction on how to make a new release for the Katib project.

## Prerequisite

1. Tools, defined in [Developer Guide](./../developer-guide.md).

1. [Write](https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/repository-permission-levels-for-an-organization#permission-levels-for-repositories-owned-by-an-organization)
   permission for the Katib repository.

1. Maintainer access to [Katib SDK](https://pypi.org/project/kubeflow-katib/).

1. Owner access to [Katib Dockerhub](https://hub.docker.com/u/kubeflowkatib).

1. Create [GitHub Token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token).

1. Install `PyGithub` to generate ChangeLog: `pip install PyGithub==1.55`

1. Install `twine` to publish SDK package: `pip install twine==3.4.1`

## Release Process

### Versioning Policy

Katib version format follows [Semantic Versioning](https://semver.org/).
Katib versions are in format of `vX.Y.Z`, where `X` is the major version, `Y` is
the minor version, and `Z` is the patch version.
The patch version contains only bug fixes.

Additionally, Katib does pre-releases in this format: `vX.Y.Z-rc.N` where N is a number
of Nth release candidate (RC) before an upcoming public release named `vX.Y.Z`.

### Release Branches and Tags

Katib releases are tagged with tags like `vX.Y.Z`, for example `v0.11.0`.

Release branches are in format of `release-X.Y`, where `X.Y` stands for
the minor release. `vX.Y.Z` releases will be released from `release-X.Y` branch.
For example, `v0.11.1` release should be on `release-0.11` branch.

If you want to push changes to the `release-X.Y` release branch, you have to
cherry pick your changes from the `master` branch and submit a PR.

### Versions for Katib Components

Katib release (git) tag includes releases for the following components:

- Manifest images with tags equal to the release
  (e.g [`v0.11.1`](https://github.com/kubeflow/katib/blob/v0.11.1/manifests/v1beta1/installs/katib-standalone/kustomization.yaml#L21-L33))

- Katib Python SDK where version is in this format: `X.Y.Z` or `X.Y.ZrcN`
  (e.g [`0.11.1`](https://github.com/kubeflow/katib/blob/v0.11.1/sdk/python/v1beta1/setup.py#L22))

### Create a new Katib Release

Follow these steps to cut a new Katib release:

1. Clone Katib repository under `$GOPATH/src` directory:

   ```
   git clone git@github.com:kubeflow/katib.git $GOPATH/src/github.com/kubeflow/katib
   ```

1. Make sure that you can build all Katib images:

   ```
   make build REGISTRY=private-registry TAG=latest
   ```

1. Create the new release:

   ```
   make release BRANCH=release-X.Y TAG=vX.Y.Z
   ```

   The above script is doing the following:

   - Create the new branch: `release-X.Y`, if it doesn't exist.

   - Create the new tag: `vX.Y.Z`.

   - Publish Katib images with the tag: `vX.Y.Z` and update manifests.

   - Publish Katib Python SDK with the version: `X.Y.Z`.

   - Push above changes to the Katib upstream `release-X.Y` branch with this commit:
     `Katib official release vX.Y.Z`

1. If the new branch was created, submit a PR to allow tests on the `release-X.Y` branch
   (e.g. [`release-0.12`](https://github.com/kubeflow/testing/pull/965))

1. Submit a PR to update SDK version on the `master` branch to the latest release.
   (e.g. [`0.12.0rc0`](TODO: ADD LINK))

1.
1. If it is not a pre-release, draft [a new GitHub Release](https://github.com/kubeflow/katib/releases/new)
   using (e.g. [Katib v0.11.0](https://github.com/kubeflow/katib/releases/tag/v0.11.0)).
