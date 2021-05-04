package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Organization struct {
	ID           githubv4.ID
	Repositories struct {
		TotalCount githubv4.Int
	} `graphql:"repositories"`
	Packages struct {
		TotalCount githubv4.Int
	} `graphql:"packages"`
	Projects struct {
		TotalCount githubv4.Int
	} `graphql:"projects"`
}

type QueryOrganization struct {
	Viewer struct {
		Organization Organization `graphql:"organization(login: $org)"`
	}
}

type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// GetOrganization retruns the organization basic information for the client
func GetOrganization(ctx context.Context, client Client) (Organization, error) {
	query := &QueryOrganization{}

	variables := map[string]interface{}{
		"org": githubv4.String("bcgov"),
	}

	if err := client.Query(ctx, query, variables); err != nil {
		fmt.Println("error:", err)
		return Organization{}, err
	}

	return query.Viewer.Organization, nil
}

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	org, _ := GetOrganization(context.Background(), client)

	fmt.Println(org)
}
