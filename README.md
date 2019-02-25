# GitHub Action to Detect an Unmergeable Pull Request

[![CircleCI](https://circleci.com/gh/cats-oss/github-action-detect-unmergeable.svg?style=svg)](https://circleci.com/gh/cats-oss/github-action-detect-unmergeable)

* This works as [GitHub Actions](https://developer.github.com/actions/).
* This detects & mark the pull request is unmergeable by changing its upstream.
* This behaves like [highfive](https://github.com/servo/highfive) or [popuko](https://github.com/voyagegroup/popuko)


## What's this?

1. This action would work if you push some changes to your repository which enables this action.
2. Then, this action would check all pull request about whether the pull request is mergeable or not.
3. If the pull request is not mergeable (unmergeable), this action do:
    * Comment to the pull request about the changeset which might breaks it. 
    * Change the label for its pull request to mark that it is unmergeable.


## Motivation

* Make easy to know what change breaks your pull request to unmergeable automatically.
* I'd like to know the changeset which breaks our pull request.


## How to use (Setup)

Add this example to your GitHub Actions workflow configuration (e.g. `.github/main.workflow`).


```
workflow "Detect unmergeable PRs" {
  on = "push"
  resolves = ["detect_unmergeable_pull_request_and_mark_them"]
}

action "detect_unmergeable_pull_request_and_mark_them" {
  uses = "cats-oss/github-action-detect-unmergeable@master"
  secrets = ["GITHUB_TOKEN"]
}
```