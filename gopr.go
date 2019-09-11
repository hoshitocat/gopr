package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var (
	token = flag.String("token", "", "Your github access token")
	head  = flag.String("head", "develop", "PR specify from `Head`. default: develop")
	base  = flag.String("base", "master", "PR specify into `Base`. default: master")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	originURL, err := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		// TODO: Error handling
		panic(err)
	}

	exp := regexp.MustCompile(`git@github\.com:(?P<owner>.+)/(?P<repo>.+)`)
	match := exp.FindSubmatch(originURL)
	originNames := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			originNames[name] = string(match[i])
		}
	}

	ownerName, repoName := originNames["owner"], strings.Split(originNames["repo"], ".git")[0]
	// TODO: テストとして `repoName` を `dotfiles` にする
	repoName = "dotfiles"
	repo, _, err := client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		// TODO: Error handling
		panic(err)
	}

	title, err := generateTitle()
	if err != nil {
		// TODO: Error handling
		panic(err)
	}
	newPR := &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(*head),
		Base:                github.String(*base),
		Body:                github.String("This is the description of the PR created with the package `github.com/google/go-github/github`"),
		MaintainerCanModify: github.Bool(true),
	}
	fmt.Println(repo)
	fmt.Println(newPR)
	//pr, _, err := client.PullRequests.Create(context.Background(), repo.Owner.GetLogin(), repo.GetName(), newPR)
	//if err != nil {
	//	// TODO: Error handling
	//	panic(err)
	//}
	//fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	comp, _, err := client.Repositories.CompareCommits(ctx, ownerName, repoName, *base, *head)
	if err != nil {
		// TODO: Error handling
		panic(err)
	}

	mergedPRMsgExp := regexp.MustCompile(`^Merge\spull\srequest\s#([0-9]+).+`)
	var mergedPRNums []int
	for _, c := range comp.Commits {
		m := mergedPRMsgExp.FindSubmatch([]byte(c.GetCommit().GetMessage()))
		if len(m) > 1 {
			n, err := strconv.Atoi(string(m[1]))
			if err != nil {
				// TODO: Error handling
				panic(err)
			}
			mergedPRNums = append(mergedPRNums, n)
		}
	}

	for _, v := range mergedPRNums {
		pr, _, err := client.PullRequests.Get(ctx, ownerName, repoName, v)
		if err != nil {
			// TODO: Error handling
			panic(err)
		}
		spew.Dump(pr)
	}
}

func generateTitle() (string, error) {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("Release %s", today), nil
}
