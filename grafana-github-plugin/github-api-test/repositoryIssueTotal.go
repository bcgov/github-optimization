package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeLineToFile(f *os.File, cells [4]string) {
	var line string
	for i, b := range cells {
		if i == 0 {
			line += b
		} else {
			line += "," + b
		}
	}
	f.WriteString(line + "\n")
}

// QueryListRepositories is the GraphQL query for retrieving a list of repositories for an organization
// {
//   search(query: "is:pr repo:grafana/grafana merged:2020-08-19..*", type: ISSUE, first: 100) {
//     nodes {
//       ... on PullRequest {
//         id
//         title
//       }
//   }
// }
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

// Repository is a code repository
type Repository struct {
	Name      string
	CreatedAt githubv4.DateTime
	UpdatedAt githubv4.DateTime
	PushedAt  githubv4.DateTime
	Issues    struct {
		TotalCount int
	}
}

// Repositories is a list of GitHub repositories
type Repositories []Repository

type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// GetRepositories retruns the organization basic information for the client
func GetRepositories(ctx context.Context, client Client) (Repositories, error) {
	var (
		variables = map[string]interface{}{
			"org":    githubv4.String("bcgov"),
			"cursor": (*githubv4.String)(nil),
		}

		repos = []Repository{}
		page  = 1
	)

	for {
		fmt.Printf("Querying %v page...\n", page)

		query := &QueryListRepositories{}
		if err := client.Query(ctx, query, variables); err != nil {
			return nil, err
		}
		r := make([]Repository, len(query.Organization.Repositories.Nodes))

		for i, v := range query.Organization.Repositories.Nodes {
			r[i] = v
		}

		repos = append(repos, r...)

		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = query.Organization.Repositories.PageInfo.EndCursor
		page++
	}

	return repos, nil
}

func main() {
	// Try creating csv first
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	err = os.MkdirAll(path+"/dat", os.ModePerm)
	check(err)

	f, err := os.Create(path + "/dat/repository-issue-total.csv")
	check(err)
	defer f.Close()

	// Main segment
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	repos, _ := GetRepositories(context.Background(), client)

	// Append data into csv
	header := [...]string{"Repository", "IssueTotalCount", "DaysOpen", "AverageIssueCountPerDay"}
	writeLineToFile(f, header)

	for _, repo := range repos {
		name := repo.Name
		issueTotalCount := repo.Issues.TotalCount
		hoursOpen := time.Now().UTC().Sub(repo.CreatedAt.UTC()).Hours()
		daysOpen := hoursOpen / 24
		averageIssueCountPerDay := float64(issueTotalCount) / daysOpen

		issueTotalCountStr := strconv.Itoa(issueTotalCount)
		daysOpenStr := strconv.Itoa(int(daysOpen))
		averageIssueCountPerDayStr := strconv.Itoa(int(averageIssueCountPerDay))

		cells := [...]string{name, issueTotalCountStr, daysOpenStr, averageIssueCountPerDayStr}

		writeLineToFile(f, cells)
	}
}
