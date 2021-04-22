package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeLineToFile(f *os.File, cells [2]string) {
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
}

// Repositories is a list of GitHub repositories
type Repositories []Repository

// RepositoryExtra is ...
type RepositoryExtra struct {
	ForkPullRequestCount int
}

// RepositoryExtras is ...
type RepositoryExtras []RepositoryExtra

type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// GetRepositories retruns the organization basic information for the client
func GetRepositories(ctx context.Context, client Client) (Repositories, RepositoryExtras, error) {
	var (
		variables = map[string]interface{}{
			"org":    githubv4.String("bcgov"),
			"cursor": (*githubv4.String)(nil),
		}

		repos  = []Repository{}
		extras = []RepositoryExtra{}
		page   = 1
	)

	for {
		fmt.Printf("Querying %v page...\n", page)

		query := &QueryListRepositories{}
		if err := client.Query(ctx, query, variables); err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
		r := make([]Repository, len(query.Organization.Repositories.Nodes))

		for i, v := range query.Organization.Repositories.Nodes {
			opts := BranchOptions{
				Org:     "bcgov",
				Name:    v.Name,
			}

			n, _ := GetForkPullRequestCount(ctx, client, opts)
			extra := RepositoryExtra{
				ForkPullRequestCount: n,
			}
			extras = append(extras, extra)

			r[i] = v
		}

		repos = append(repos, r...)

		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = query.Organization.Repositories.PageInfo.EndCursor
		page++

		// time.Sleep(10 * time.Minute)
	}

	return repos, extras, nil
}

type BranchOptions struct {
	Org     string
	Name    string
}

// QueryForkPullRequestCount is ...
type QueryForkPullRequestCount struct {
	Repository struct {
		PullRequests struct {
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
			Nodes []struct {
				Repository struct {
					Name string
				}
				BaseRepository struct {
					Name string
				}
				HeadRepository struct {
					Name string
				}
			}
		} `graphql:"pullRequests(first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $org)"`
}

// GetForkPullRequestCount retruns ...
func GetForkPullRequestCount(ctx context.Context, client Client, ops BranchOptions) (int, error) {
	var (
		variables = map[string]interface{}{
			"org":     githubv4.String(ops.Org),
			"name":    githubv4.String(ops.Name),
			"cursor": (*githubv4.String)(nil),
		}

		forkPrCount  = 0
		page   = 1
	)

	for {
		fmt.Printf("Querying %v - %v pr page...\n", ops.Name, page)

		query := &QueryForkPullRequestCount{}
		if err := client.Query(ctx, query, variables); err != nil {
			fmt.Println(err)
			return 0, err
		}

		for _, v := range query.Repository.PullRequests.Nodes {
			if v.BaseRepository.Name != v.HeadRepository.Name {
				forkPrCount++
			}
		}

		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = query.Repository.PullRequests.PageInfo.EndCursor
		page++
	}

	return forkPrCount,  nil
}

func main() {
	// Try creating csv first
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	targetDir := "../../../notebook/dat/"
	targetFile := "/repository-fork-pr.csv"

	err = os.MkdirAll(path+targetDir, os.ModePerm)
	check(err)

	f, err := os.Create(path + targetDir + targetFile)
	check(err)
	defer f.Close()

	// Main segment
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	repos, extras, _ := GetRepositories(context.Background(), client)

	// Append data into csv
	header := [...]string{
		"Repository",
		"Fork PullRequest Count",
	}
	writeLineToFile(f, header)

	for i, repo := range repos {
		name := repo.Name
		forkPrCount := extras[i].ForkPullRequestCount
		forkPrCountStr := strconv.Itoa(forkPrCount)

		cells := [...]string{
			name,
			forkPrCountStr,
		}

		writeLineToFile(f, cells)
	}
}
