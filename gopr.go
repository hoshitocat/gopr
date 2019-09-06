package main

import (
	"context"
	"flag"
	"os/exec"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

var (
	token = flag.String("token", "", "Your github access token")
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
	repo, _, err := client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		// TODO: Error handling
		panic(err)
	}

	spew.Dump(repo)
}
