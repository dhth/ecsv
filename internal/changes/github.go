package changes

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/dhth/ecsv/internal/types"
	"github.com/google/go-github/v72/github"
)

func FetchChanges(
	client *github.Client,
	systemKey,
	owner,
	repo,
	baseRef,
	headRef string,
	ignorePattern *regexp.Regexp,
) types.ChangesResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	options := github.ListOptions{
		Page:    0,
		PerPage: 100,
	}

	comparison, _, err := client.Repositories.CompareCommits(ctx, owner, repo, "v"+baseRef, "v"+headRef, &options)
	if err != nil {
		return types.ChangesResult{
			SystemKey: systemKey,
			Error:     err,
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

		if ignorePattern != nil && ignorePattern.Match([]byte(commit.Commit.GetMessage())) {
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
		SystemKey: systemKey,
		Commits:   commits,
	}
}
