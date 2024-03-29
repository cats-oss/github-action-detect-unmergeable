# GitHub Action to Detect an Unmergeable Pull Request

![CI](https://github.com/cats-oss/github-action-detect-unmergeable/workflows/CI/badge.svg)

* This works as [GitHub Actions](https://help.github.com/en/articles/about-github-actions).
* This detects & mark the pull request is unmergeable by changing its upstream.
* This behaves like [highfive](https://github.com/servo/highfive) or [popuko](https://github.com/voyagegroup/popuko)


## What's this?

1. This action would work if you push some changes to your repository which enables this action.
2. Then, this action would check all pull request about whether the pull request is mergeable or not.
3. If the pull request is not mergeable (unmergeable), this action do:
    * Comment to the pull request about the changeset which might breaks it. 
    * Change the label for its pull request to mark that it is unmergeable.
        * This action can remove the added label if the PR's conflict is resolved by rebasing or others (_optional_).


## Motivation

* Make easy to know what change breaks your pull request to unmergeable automatically.
* I'd like to know the changeset which breaks our pull request.


## How to use (Setup)

Add this example to your GitHub Actions [workflow configuration](https://help.github.com/en/articles/configuring-workflows).


### YAML syntax

```yaml
name: Detect unmergeable PRs

on:
  push:
    branches:
      - "*"
    # Ignore all pushing for tags
    tags:
      - "!*"
  # If you'd like to remove the added label by this action automatically
  # on updating a pull request by pushing changes.
  # Please recieve this event.
  pull_request:
    types: synchronize

jobs:
  detect_unmergeable_pull_request_and_mark_them:
    runs-on: ubuntu-latest
    permissions:
      contents: none
      pull-requests: write
    steps:
      - name: Run the action to detect unmergeable PRs
        # We recommend to use an arbitary latest version
        # if you don't have any troubles.
        # You can also specify `master`, but it sometimes might be broken.
        uses: cats-oss/github-action-detect-unmergeable@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Debugging

If you have some troubles, please try to see information
by inserting the below snippet to `steps` section for this workflow.

```yaml
      - name: Dump GitHub Context
        env:
          GITHUB_CONTEXT: ${{ toJson(github) }}
        run: echo "${GITHUB_CONTEXT}"
```