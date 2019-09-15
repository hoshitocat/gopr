package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
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

var client *github.Client

func init() {
}

func main() {
	flag.Parse()
	ctx := context.Background()

	client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})))

	owner, repo, err := getCurrentRepo(ctx)
	if err != nil {
		fmt.Printf("gopr: %+v\n", err.Error())
		os.Exit(1)
	}

	title, err := generateTitle()
	if err != nil {
		fmt.Printf("gopr: %+v\n", err.Error())
		os.Exit(1)
	}

	body, err := generateBody(ctx, owner, repo)
	if err != nil {
		fmt.Printf("gopr: %+v\n", err.Error())
		os.Exit(1)
	}

	pr, _, err := client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title:               github.String(title),
		Head:                github.String(*head),
		Base:                github.String(*base),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	})
	if err != nil {
		fmt.Printf("gopr: %+v\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
}

func getCurrentRepo(ctx context.Context) (owner string, repo string, err error) {
	originURL, err := exec.CommandContext(ctx, "git", "config", "--get", "remote.origin.url").Output()
	if err != nil || len(originURL) == 0 {
		return "", "", errors.New("could not get current git repository origin URL")
	}

	exp := regexp.MustCompile(`git@github\.com:(?P<owner>.+)/(?P<repo>.+)`)
	match := exp.FindSubmatch(originURL)
	originNames := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			originNames[name] = string(match[i])
		}
	}

	owner, repo = originNames["owner"], strings.Split(originNames["repo"], ".git")[0]
	// TODO: テストとして `repoName` を `dotfiles` にする
	repo = "dotfiles"
	return owner, repo, nil
}

func generateBody(ctx context.Context, owner, repo string) (string, error) {
	comp, _, err := client.Repositories.CompareCommits(ctx, owner, repo, *base, *head)
	if err != nil {
		return "", errors.New(fmt.Sprintf("could not genereate PR body: %+v", err))
	}

	mergedPRMsgExp := regexp.MustCompile(`^Merge\spull\srequest\s#([0-9]+).+`)
	var mergedPRNums []int
	for _, c := range comp.Commits {
		m := mergedPRMsgExp.FindSubmatch([]byte(c.GetCommit().GetMessage()))
		if len(m) > 1 {
			n, err := strconv.Atoi(string(m[1]))
			if err != nil {
				return "", errors.New(fmt.Sprintf("could not genereate PR body: %+v", err))
			}
			mergedPRNums = append(mergedPRNums, n)
		}
	}

	var body string
	for _, v := range mergedPRNums {
		pr, _, err := client.PullRequests.Get(ctx, owner, repo, v)
		if err != nil {
			return "", errors.New(fmt.Sprintf("could not genereate PR body: %+v", err))
		}
		body += fmt.Sprintf("- [ ] [#%d](%s) %s created by @%s\n", v, pr.GetHTMLURL(), pr.GetTitle(), pr.GetUser().GetLogin())
	}

	return body, nil
}

func generateTitle() (string, error) {
	today := time.Now().Format("2006-01-02")

	content, err := ioutil.ReadFile(".github/TITLE_TEMPLATE.md")
	if err != nil {
		return "", errors.New(fmt.Sprintf("could not generate PR title: %+v", err))
	}
	return fmt.Sprintf(string(content), today), nil
}
