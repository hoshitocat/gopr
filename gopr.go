package main

import (
	"context"
	"flag"
	"fmt"
	"os/exec"

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

	_, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		// TODO: Error handling
		panic(err)
	}

	originURL, err := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		// TODO: Error handling
		panic(err)
	}
	repoName := string(originURL)

	fmt.Println(repoName)
}
