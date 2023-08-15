# Release the Katib Project

This is the instruction on how to make a new release for the Katib project.

## Prerequisite

- Tools, defined in the [Developer Guide](./../developer-guide.md#requirements).

- [Write](https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/repository-permission-levels-for-an-organization#permission-levels-for-repositories-owned-by-an-organization)
  permission for the Katib repository.

- Maintainer access to the [Katib SDK](https://pypi.org/project/kubeflow-katib/).

- Owner access to the [Katib Dockerhub](https://hub.docker.com/u/kubeflowkatib).

- Create a [GitHub Token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token).

- Install `PyGithub` to generate the [Changelog](./../../CHANGELOG.md): `pip install PyGithub==1.55`

- Install `twine` to publish the SDK package: `pip install twine==3.4.1`

  - Create a [PyPI Token](https://pypi.org/help/#apitoken) to publish Katib SDK.

  - Add the following config to your `~/.pypirc` file:

    ```
    [pypi]
       username = __token__
       password = <PYPI_TOKEN>
    ```

## Release Process

### Versioning Policy

Katib version format follows [Semantic Versioning](https://semver.org/).
Katib versions are in the format of `vX.Y.Z`, where `X` is the major version, `Y` is
the minor version, and `Z` is the patch version.
The patch version contains only bug fixes.

Additionally, Katib does pre-releases in this format: `vX.Y.Z-rc.N` where `N` is a number
of the `Nth` release candidate (RC) before an upcoming public release named `vX.Y.Z`.

### Release Branches and Tags

Katib releases are tagged with tags like `vX.Y.Z`, for example `v0.11.0`.

Release branches are in the format of `release-X.Y`, where `X.Y` stands for
the minor release.

`vX.Y.Z` releases are released from the `release-X.Y` branch. For example,
`v0.11.1` release should be on `release-0.11` branch.

If you want to push changes to the `release-X.Y` release branch, you have to
cherry pick your changes from the `master` branch and submit a PR.

### Versions for Katib Components

Katib release ([git tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging))
includes releases for the following components:

- Manifest images with tags equal to the release
  (e.g [`v0.11.1`](https://github.com/kubeflow/katib/blob/v0.11.1/manifests/v1beta1/installs/katib-standalone/kustomization.yaml#L21-L33)).

- Katib Python SDK where version is in this format: `X.Y.Z` or `X.Y.ZrcN`
  (e.g [`0.11.1`](https://github.com/kubeflow/katib/blob/v0.11.1/sdk/python/v1beta1/setup.py#L22)).

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

   - Create the new tag: `vX.Y.Z` from the release branch: `release-X.Y`.

   - Publish Katib images with the tag: `vX.Y.Z` and update manifests.

   - Publish Katib Python SDK with the version: `X.Y.Z`.

   - Push above changes to the Katib upstream `release-X.Y` branch with this commit:
     `Katib official release vX.Y.Z`

1. Submit a PR to update the SDK version on the `master` branch to the latest release.
   (e.g. [`#1640`](https://github.com/kubeflow/katib/pull/1640)).

1. Update the Changelog by running:

   ```
   python docs/release/changelog.py --token=<github-token> --range=<previous-release>..<current-release>
   ```

   If you are creating the **first minor pre-release** or the **minor** release (`X.Y`), your
   `previous-release` is equal to the latest release on the `release-X.Y-1` branch.
   For example: `--range=v0.11.1..v0.12.0`

   Otherwise, your `previous-release` is equal to the latest release on the `release-X.Y` branch.
   For example: `--range=v0.12.0-rc.0..v0.12.0-rc.1`

   Group PRs in the Changelog into Features, Bug fixes, Documentation, etc.
   Check this example: [v0.11.0](https://github.com/kubeflow/katib/releases/tag/v0.11.0)

   Finally, submit a PR with the updated Changelog.

1. If it is not a pre-release, draft [a new GitHub Release](https://github.com/kubeflow/katib/releases/new).
