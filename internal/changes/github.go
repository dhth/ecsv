package changes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dhth/ecsv/internal/types"
	"github.com/google/go-github/v72/github"
)

var errCouldntGetTokenFromGH = errors.New("couldn't get token from GitHub's CLI")

const transformPlaceholder = "{{version}}"

func GetGHClient() (*github.Client, error) {
	var zero *github.Client
	tokenFromEnv := os.Getenv("GH_TOKEN")
	if tokenFromEnv != "" {
		return github.NewClient(nil).WithAuthToken(os.Getenv("GH_TOKEN")), nil
	}

	tokenFromGH, err := getTokenFromGH()
	if err != nil {
		return zero, err
	}

	client := github.NewClient(nil).WithAuthToken(tokenFromGH)
	return client, nil
}

func getTokenFromGH() (string, error) {
	var zero string
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntGetTokenFromGH, err.Error())
	}

	return strings.TrimSpace(string(output)), nil
}

func FetchChanges(
	client *github.Client,
	config types.ChangesConfig,
	baseRef,
	headRef string,
) types.ChangesResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := github.ListOptions{
		Page:    0,
		PerPage: 100,
	}

	baseRefToUse := baseRef
	if config.Transform != nil {
		baseRefToUse = strings.Replace(*config.Transform, transformPlaceholder, baseRef, 1)
	}

	headRefToUse := headRef
	if config.Transform != nil {
		headRefToUse = strings.Replace(*config.Transform, transformPlaceholder, headRef, 1)
	}

	comparison, _, err := client.Repositories.CompareCommits(ctx, config.Owner, config.Repo, baseRefToUse, headRefToUse, &options)
	if err != nil {
		return types.ChangesResult{
			Config: config,
			Error:  err,
		}
	}

	//nolint:prealloc
	var commits []types.Commit
	for _, commit := range comparison.Commits {
		author := commit.GetCommit().GetAuthor()
		var ca string
		var at string

		if author != nil {
			ca = author.GetName()
			at = author.GetDate().Format(time.RFC3339)
		}
		sha := commit.GetSHA()
		if len(sha) > 8 {
			sha = sha[:8]
		}

		if config.IgnorePattern != nil && config.IgnorePattern.Match([]byte(commit.Commit.GetMessage())) {
			continue
		}

		message := strings.Split(commit.Commit.GetMessage(), "\n")[0]

		commits = append(commits, types.Commit{
			SHA:        sha,
			Message:    message,
			HTMLURL:    commit.GetHTMLURL(),
			Author:     ca,
			AuthoredAt: at,
		})

	}

	return types.ChangesResult{
		Config:  config,
		Commits: commits,
		DiffURL: fmt.Sprintf("https://github.com/%s/%s/compare/%s...%s", config.Owner, config.Repo, baseRefToUse, headRefToUse),
	}
}
