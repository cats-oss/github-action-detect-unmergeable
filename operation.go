package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v24/github"
)

func getOpenPullRequestAll(client *github.Client, owner, name string) []*github.PullRequest {
	ctx := context.Background()
	list, _, err := client.PullRequests.List(ctx, owner, name, &github.PullRequestListOptions{
		State: "open",
	})

	if err != nil {
		log.Printf("%v", err)
		return nil
	}

	return list
}

func isRelatedToPushedBranch(pullReqInfo *github.PullRequest, pushedBranchRef string) (ok bool, hasRelationShip bool) {
	//  Define to explain followings comments:
	//      * _origin_ is the repository which registers & runs this action.
	//      * _forked_ is the repository which is forked from _origin_.
	//
	//  We don't have to imagine the case which the user pushed to the branch in _forked_.
	//  and its branch is opened as the pull requet for _origin_.
	//  Because _push_ event happens and this action will run only if someone pushed to a branch on _origin_.
	//
	//  By the observation of the behavior of GitHub REST API v3 Pull Requests,
	//  object's values would be followings if `repository_name:branch_name` is `octocat:master`.
	//
	//      * `.base.ref`
	//          * The value is `master`.
	//          * This value is the branch name which the pull request would be merged into.
	//            even if the pull request is came from _forked_ or the one's parent is arbitary commit.
	//      * `.base.label` is `octocat:master`.
	//          * The value is `octocat:master`.
	//          * This value is same even if the pull request is came from _forked_ or the one's parent is arbitary commit.
	//
	//  By theese result, I think we can judge that the pull request is related to the pushed branch by checking
	//  whether `.base.ref` is same with the name of the pushed branch.

	base := pullReqInfo.GetBase()
	if base == nil {
		ok = false
		return
	}

	currentBranchName := base.GetRef()
	if currentBranchName == "" {
		ok = false
		log.Println("currentBranchName is empty")
		return
	}

	pushedBranchName := strings.Replace(pushedBranchRef, "refs/heads/", "", -1)
	if pushedBranchName == "" {
		ok = false
		log.Println("pushedBranchName is empty")
		return
	}

	ok = true
	hasRelationShip = currentBranchName == pushedBranchName
	return
}

func checkAndMarkIfPullRequestUnmergeable(client *github.Client, owner, repo string, oldPR *github.PullRequest, compareURL string) {
	number := oldPR.GetNumber()
	if number == 0 {
		return
	}

	if hasNeedRebaseLabel(oldPR.Labels) {
		log.Printf("#%v has been labeled as %v\n", number, labelStatusNeedRebase)
		return
	}

	// When we get all opened pull requests, GitHub have not checked whether ths PR is unmergeble yet.
	// So I think we should retry them.
	hasCompleted, shouldMark, err := shouldMarkPullRequestNeedRebase(client, owner, repo, number)
	if err != nil {
		log.Printf("#%v fails: %v\n", number, err)
		return
	}

	if !hasCompleted {
		// retry
		time.Sleep(10 * time.Second)
		hasCompleted, shouldMark, err = shouldMarkPullRequestNeedRebase(client, owner, repo, number)
		if err != nil {
			log.Printf("#%v fails: %v\n", number, err)
			return
		}
	}

	log.Printf("#%v is hasCompleted: %v, shouldMark: %v\n", number, hasCompleted, shouldMark)
	if !hasCompleted {
		// give up
		return
	}

	if !shouldMark {
		return
	}

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		labels := []string{labelStatusNeedRebase}
		_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, number, labels)
		if err != nil {
			log.Printf("#%v: Fail to add labels: %v\n", number, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		body := fmt.Sprintf(`:umbrella: The latest upstream change (presumably [these](%v)) made this pull request unmergeable. Please resolve the merge conflicts.`, compareURL)
		_, _, err := client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{
			Body: &body,
		})

		if err != nil {
			log.Printf("#%v: Fail to add comments: %v\n", number, err)
		}
	}()

	wg.Wait()
}

const labelStatusNeedRebase = "S-needs-rebase"

func hasNeedRebaseLabel(labels []*github.Label) bool {
	for _, label := range labels {
		name := label.GetName()
		if name == labelStatusNeedRebase {
			return true
		}
	}

	return false
}

func shouldMarkPullRequestNeedRebase(client *github.Client, owner, repo string, number int) (hasCompleted bool, shouldMark bool, err error) {
	ctx := context.Background()
	newPRInfoResponse, _, err := client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return
	}

	mergeable := newPRInfoResponse.Mergeable
	if mergeable == nil {
		return
	}

	hasCompleted = true

	// Check again to confirm the other instance of this action's behavior.
	if hasNeedRebaseLabel(newPRInfoResponse.Labels) {
		shouldMark = false
		return
	}

	if *mergeable {
		shouldMark = false
	} else {
		shouldMark = true
	}

	return
}
