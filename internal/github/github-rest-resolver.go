package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// rest api
func GetContributors(repoName string, orgName string, token string) ([]*github.ContributorStats, error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := github.NewClient(httpClient)

	contributors, _, err := client.Repositories.ListContributorsStats(ctx, orgName, repoName)
	return contributors, err
}
