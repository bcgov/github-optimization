package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type QueryListRepositories struct {
	Organization struct {
		Repositories struct {
			TotalCount githubv4.Int
		}
	} `graphql:"organization(login: $org)"`
}

// type QueryListRepositories struct {
// 	Organization struct {
// 		Repositories []struct {
// 			PageInfo struct {
// 				HasNextPage bool
// 				EndCursor   githubv4.String
// 			}
// 			nodes struct {
// 				Name githubv4.String
// 			}
// 		} `graphql:"repositories(first: 100)"`
// 	} `graphql:"organization(login: $org)"`
// }

func main() {
	q := &QueryListRepositories{}

	variables := map[string]interface{}{
		"org": githubv4.String("bcgov"),
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
	fmt.Print(q.Organization.Repositories.TotalCount)
}
