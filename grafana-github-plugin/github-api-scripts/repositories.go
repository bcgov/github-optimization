package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Repository struct {
	Name  string
	Owner struct {
		Login string
	}
	NameWithOwner      string
	URL                string
	HomepageURL        string
	Description        string
	ForkCount          int64
	IsFork             bool
	IsMirror           bool
	IsPrivate          bool
	IsArchived         bool
	IsTemplate         bool
	StargazerCount     int64
	DiskUsage          int64
	HasIssuesEnabled   bool
	HasProjectsEnabled bool
	HasWikiEnabled     bool
	MergeCommitAllowed bool
	RebaseMergeAllowed bool
	SquashMergeAllowed bool
	CreatedAt          githubv4.DateTime
	UpdatedAt          githubv4.DateTime
	PushedAt           githubv4.DateTime
}

type QueryListRepositories struct {
	Organization struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
			Nodes []Repository
		} `graphql:"repositories(first: 100, after: $cursor)"`
	} `graphql:"organization(login: $org)"`
}

func main() {
	q := &QueryListRepositories{}

	variables := map[string]interface{}{
		"org":    githubv4.String("bcgov"),
		"cursor": (*githubv4.String)(nil),
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	fmt.Print(client)
	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		// Handle error.
		fmt.Print(err)
	}
	for i, b := range q.Organization.Repositories.Nodes {
		description := b.Description
		fmt.Println("description", i, description)
		fmt.Println("name", i, b.NameWithOwner)
		fmt.Println("page", i, b.HomepageURL)
		fmt.Println("is spoon", i, b.IsFork)
	}

}
