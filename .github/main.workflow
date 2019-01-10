workflow "Detect unmergeable PRs" {
  on = "push"
  resolves = ["detect_unmergeable_pull_request_and_mark_them"]
}

action "detect_unmergeable_pull_request_and_mark_them" {
  uses = "cats-oss/github-action-detect-unmergeable@application-code"
  secrets = ["GITHUB_TOKEN"]
  env = {
    ENABLE_CHECK_RELATIONSHIP_BETWEEN_PUSHED_AND_PR = "true"
  }
}