# Releasing Kubeflow Katib

## Prerequisites

- [Write](https://docs.github.com/en/organizations/managing-access-to-your-organizations-repositories/repository-permission-levels-for-an-organization#permission-levels-for-repositories-owned-by-an-organization)
  permission for the Katib repository.

- GitHub **`release` environment** with required reviewers (gates PyPI and GitHub Release jobs).

- [PyPI trusted publishing](https://docs.pypi.org/trusted-publishers/) for `kubeflow-katib` and `kubeflow_katib_api`
  (workflow: `release.yaml`, owner: `kubeflow`, repo: `katib`).

- Repository secrets: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`.

- Optional: [GitHub token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)
  as `GITHUB_TOKEN` and `git-cliff` for changelog generation via `make release`.

## Versioning Policy

Katib follows [Semantic Versioning](https://semver.org/) and Python [PEP 440](https://peps.python.org/pep-0440/) for SDK packages.

| Artifact | Format | Example |
| --- | --- | --- |
| Git tag | `vX.Y.Z` or `vX.Y.Z-rc.N` | `v0.19.1`, `v0.19.0-rc.0` |
| Python SDK / API | `X.Y.Z` or `X.Y.ZrcN` | `0.19.1`, `0.19.0rc0` |
| Release branch | `release-X.Y` | `release-0.19` |

Pre-releases use Python `X.Y.ZrcN` and git tags `vX.Y.Z-rc.N`.

## Release Branches

Release branches use the format `release-X.Y` (for example `release-0.19`).

- **Latest minor series**: bump the version on `master` and open a PR to `master`.
- **Older minor patch** (for example `0.18.1` while `master` is at `0.19.x`): open a PR to the
  corresponding `release-0.18` branch (backport fixes via PRs, not manual cherry-picks).

Manifest image tags stay at `latest` on `master`. CI pins them to the release tag on the
`release-X.Y` branch during the automated release workflow.

## Step-by-Step Release Process

### 1. Update version and changelog

```sh
export GITHUB_TOKEN=<your_github_token>   # optional, for git-cliff rate limits
make release VERSION=<X.Y.Z>
# e.g. make release VERSION=0.19.1
```

This updates:

- `sdk/python/v1beta1/setup.py`
- `hack/python-api/gen-api.sh` and `api/python_api/kubeflow_katib_api/__init__.py`
- `CHANGELOG.md` (stable releases only, when `git-cliff` is installed)

Review the diff (`git diff`) and open a PR:

- **Latest minor series** → PR to `master`
- **Older patch** → PR to `release-X.Y`

Wait for the [Check Release](https://github.com/kubeflow/katib/actions/workflows/check-release.yaml) workflow.

Before merge, run the [Release workflow](https://github.com/kubeflow/katib/actions/workflows/release.yaml) on your PR branch with **`dry_run: true`** (default). This runs the same Prepare and Build jobs as a real release without pushing branches, tags, or publishing artifacts.

### 2. Automated release

Merge the PR. A push that changes `sdk/python/v1beta1/setup.py` triggers the Release workflow, which:

1. **Prepare** — creates or updates `release-X.Y`, pins manifest image tags on that branch
2. **Build** — validates versions, builds Python packages
3. **Tag** — creates and pushes `vX.Y.Z`
4. **Publish images** — multi-arch images to GHCR and DockerHub

Confirm the release branch and tag appear on GitHub.

### 3. Manual approvals

1. [GitHub Actions](https://github.com/kubeflow/katib/actions) → **Release** → approve **Publish to PyPI**
2. After PyPI succeeds → approve **Create GitHub Release**

### 4. Verify

- [GHCR](https://github.com/kubeflow/katib/pkgs/container/katib) and [DockerHub](https://hub.docker.com/u/kubeflowkatib) images
- [PyPI kubeflow-katib](https://pypi.org/project/kubeflow-katib/) and [kubeflow-katib-api](https://pypi.org/project/kubeflow-katib-api/)
- [GitHub Releases](https://github.com/kubeflow/katib/releases)
- `pip install kubeflow-katib==X.Y.Z`

## Announcement

For minor/major releases, announce on Kubeflow community channels
([Slack](https://www.kubeflow.org/docs/about/community/#kubeflow-slack-channels),
[mailing list](https://www.kubeflow.org/docs/about/community/#kubeflow-mailing-list)).
