package main

import (
	"context"
	"flag"
	"fmt"

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
	repos, _, err := client.Repositories.List(ctx, "", nil)
	if err != nil {
		// TODO: Error handling
		panic(err)
	}

	fmt.Println(repos)
}
