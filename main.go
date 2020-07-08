package main

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

const defaultNeedRebaseLabel = "S-needs-rebase"

const (
	ACTION_TYPE_PUSH     = iota
	ACTION_TYPE_PULL_REQ = iota
)

func main() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatalln("$GITHUB_TOKEN is empty")
		return
	}

	repositoryFullName := os.Getenv("GITHUB_REPOSITORY")
	if repositoryFullName == "" {
		log.Fatalln("$GITHUB_REPOSITORY is empty")
		return
	}

	var actionType int
	githubEventName := os.Getenv("GITHUB_EVENT_NAME")
	switch githubEventName {
	case "":
		log.Fatalln("$GITHUB_EVENT_NAME is empty")
		return
	case "push":
		actionType = ACTION_TYPE_PUSH
	case "pull_request":
		actionType = ACTION_TYPE_PULL_REQ
	default:
		log.Fatalln("Unsupported $GITHUB_EVENT_NAME: " + githubEventName)
		return
	}

	githubEventPath := os.Getenv("GITHUB_EVENT_PATH")
	if githubEventPath == "" {
		log.Fatalln("$GITHUB_EVENT_PATH is empty")
		return
	}

	needRebaseLabel := os.Getenv("LABEL_NEED_REBASE")
	if needRebaseLabel == "" {
		needRebaseLabel = defaultNeedRebaseLabel
	}
	log.Printf("We will supply `%v` if the pull request is unmergeable\n", needRebaseLabel)

	githubClient := createGithubClient(githubToken)
	if githubClient == nil {
		log.Fatalln("could not create githubClient")
		return
	}

	repoPair := strings.Split(repositoryFullName, "/")
	repoOwner := repoPair[0]
	repoName := repoPair[1]

	switch actionType {
	case ACTION_TYPE_PUSH:
		log.Println("Search and mark unmergeable pull requests.")
		onPushEvent(githubClient, githubEventPath, repoOwner, repoName, needRebaseLabel)
	case ACTION_TYPE_PULL_REQ:
		log.Println("Check whether the synced pull request is mergeable or not.")
		onPullRequestEvent(githubClient, githubEventPath, repoOwner, repoName, needRebaseLabel)
	default:
		return
	}
}

func createGithubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}

func onPushEvent(githubClient *github.Client, githubEventPath string, repoOwner string, repoName string, needRebaseLabel string) {
	eventData := loadJSONFileForPushEventData(githubEventPath)
	if eventData == nil {
		log.Fatal("Could not get eventData")
		return
	}

	eventOriginRefName := eventData.GetRef()
	if eventOriginRefName == "" {
		log.Println("eventOriginRefName is empty string")
		return
	}

	compareURL := eventData.GetCompare()
	if compareURL == "" {
		// avoid this.
		log.Println("compareURL is empty string")
	}

	openedPRList := getOpenPullRequestAll(githubClient, repoOwner, repoName)
	if openedPRList == nil {
		return
	}

	wg := &sync.WaitGroup{}
	for _, pr := range openedPRList {
		wg.Add(1)
		go func(pr *github.PullRequest) {
			defer wg.Done()

			number := pr.GetNumber()
			ok, hasRelationShip := isRelatedToPushedBranch(pr, eventOriginRefName)
			if !ok {
				log.Printf("#%v mysterious result and abort it \n", number)
				return
			}

			if !hasRelationShip {
				log.Printf("#%v is not related to %v \n", number, eventOriginRefName)
				return
			}

			checkAndMarkIfPullRequestUnmergeable(githubClient, repoOwner, repoName, pr, compareURL, needRebaseLabel)
		}(pr)
	}
	wg.Wait()
}

func onPullRequestEvent(githubClient *github.Client, githubEventPath string, repoOwner string, repoName string, needRebaseLabel string) {
	eventData := loadJSONFileForPullRequestEventData(githubEventPath)
	if eventData == nil {
		log.Fatal("Could not get eventData")
		return
	}

	if action := eventData.GetAction(); action != "synchronize" {
		log.Printf("we cannot handle `#%v`", action)
		return
	}

	prNumber := eventData.GetNumber()
	if prNumber == 0 {
		log.Println("we cannot get the PR number for this pull request")
		return
	}

	_, isUnmergeable, err := shouldMarkPullRequestNeedRebase(githubClient, repoOwner, repoName, prNumber, needRebaseLabel)
	if err != nil {
		return
	}

	if isUnmergeable {
		log.Printf("#%v is not mergeable. We don't do anything", prNumber)
		return
	}

	ctx := context.Background()
	if _, err := githubClient.Issues.RemoveLabelForIssue(ctx, repoOwner, repoName, prNumber, needRebaseLabel); err != nil {
		log.Printf("#%v is mergeable but fail to remove the label `%v`", prNumber, needRebaseLabel)
		return
	}

	log.Printf("#%v is mergeable. We removed the label `%v`", prNumber, needRebaseLabel)
}
