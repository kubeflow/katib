from github import Github
import argparse

REPO_NAME = "kubeflow/katib"

parser = argparse.ArgumentParser()
parser.add_argument("--token", type=str, help="GitHub Access Token")
parser.add_argument("--range", type=str, help="Changelog is generated for this release range")
args = parser.parse_args()

if args.token is None:
    raise Exception("GitHub Token must be set")

try:
    changes_from = args.range.split("..")[0]
    change_to = args.range.split("..")[1]
except Exception:
    raise Exception("Release range must be set in this format: v0.11.0..v0.12.0")


github_repo = Github(args.token).get_repo(REPO_NAME)
commits = github_repo.compare(changes_from, change_to).commits

# for commit in commits:
#     # Only add commits with PRs.
#     for pr in commit.get_pulls():
#         print(pr.title)
#         print(pr.id)
#         print(pr.html_url)


f = open("./../CHANGELOG.md", "a")
f.write("Data")
f.close()
