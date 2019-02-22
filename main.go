package main

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
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

	githubEventName := os.Getenv("GITHUB_EVENT_NAME")
	if githubEventName == "" {
		log.Fatalln("$GITHUB_EVENT_NAME is empty")
		return
	}
	if githubEventName != "push" {
		log.Fatalln("Unsupported $GITHUB_EVENT_NAME: " + githubEventName)
		return
	}

	githubEventPath := os.Getenv("GITHUB_EVENT_PATH")
	if githubEventPath == "" {
		log.Fatalln("$GITHUB_EVENT_PATH is empty")
		return
	}

	eventData := loadJSONFile(githubEventPath)
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

	githubClient := createGithubClient(githubToken)
	if githubClient == nil {
		log.Fatalln("could not create githubClient")
		return
	}

	repoPair := strings.Split(repositoryFullName, "/")
	repoOwner := repoPair[0]
	repoName := repoPair[1]

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

			checkAndMarkIfPullRequestUnmergeable(githubClient, repoOwner, repoName, pr, compareURL)
		}(pr)
	}
	wg.Wait()
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
