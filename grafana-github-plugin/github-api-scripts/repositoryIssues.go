package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	utils "github.com/grafana/github-datasource/github-api-scripts/utils"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

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

// Repository is ...
type Repository struct {
	Name string
}

// Repo is ...
type Repo struct {
	Name   string
	Issues Issues
}

// Repositories is a list of GitHub repositories
type Repositories []Repo

// Issue is ...
type Issue struct {
	Author struct {
		Login string
	}
	Closed       bool
	ClosedAt     string
	CreatedAt    string
	LastEditedAt string
	State        string
	Title        string
	UpdatedAt    string
}

// Issues is ...
type Issues []Issue

type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

type RepositoryOptions struct {
	Org string
}

// GetRepositories retruns the organization basic information for the client
func GetRepositories(ctx context.Context, client Client, opts RepositoryOptions) (Repositories, error) {
	var (
		variables = map[string]interface{}{
			"org":    githubv4.String(opts.Org),
			"cursor": (*githubv4.String)(nil),
		}

		repos = []Repo{}
		page  = 1
	)

	for {
		fmt.Printf("Querying %v page...\n", page)

		query := &QueryListRepositories{}
		if err := client.Query(ctx, query, variables); err != nil {
			fmt.Println(err)
			return []Repo{}, err
		}

		r := make([]Repo, len(query.Organization.Repositories.Nodes))

		for i, v := range query.Organization.Repositories.Nodes {
			opts := BranchOptions{
				Org:  opts.Org,
				Name: v.Name,
			}

			issues, _ := GetRepositoryIssues(ctx, client, opts)

			repo := Repo{}
			repo.Name = v.Name
			repo.Issues = issues

			r[i] = repo
		}

		repos = append(repos, r...)

		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = query.Organization.Repositories.PageInfo.EndCursor
		page++

		// time.Sleep(10 * time.Minute)
	}

	return repos, nil
}

type BranchOptions struct {
	Org  string
	Name string
}

// QueryRepositoryIssues is ...
type QueryRepositoryIssues struct {
	Repository struct {
		Issues struct {
			Nodes    []Issue
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
		} `graphql:"issues(first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $org)"`
}

// GetRepositoryIssues retruns ...
func GetRepositoryIssues(ctx context.Context, client Client, ops BranchOptions) (Issues, error) {
	var (
		variables = map[string]interface{}{
			"org":    githubv4.String(ops.Org),
			"name":   githubv4.String(ops.Name),
			"cursor": (*githubv4.String)(nil),
		}

		issues = []Issue{}

		page = 1
	)

	for {
		fmt.Printf("Querying %v - %v issue page...\n", ops.Name, page)

		query := &QueryRepositoryIssues{}
		if err := client.Query(ctx, query, variables); err != nil {
			fmt.Println(err)
			return []Issue{}, err
		}

		issues = append(issues, query.Repository.Issues.Nodes...)

		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = query.Repository.Issues.PageInfo.EndCursor
		page++
	}

	return issues, nil
}

func main() {
	token, org := utils.CheckEnv()

	// Try creating csv first
	path, err := os.Getwd()
	utils.HandleError(err)

	targetDir := fmt.Sprintf("../../../notebook/dat/%v/", org)
	targetFile := "/repository-issues.csv"

	err = os.MkdirAll(path+targetDir, os.ModePerm)
	utils.HandleError(err)

	f, err := os.Create(path + targetDir + targetFile)
	utils.HandleError(err)
	defer f.Close()

	// Main segment
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	opts := RepositoryOptions{
		Org: org,
	}

	repos, _ := GetRepositories(context.Background(), client, opts)

	// Append data into csv
	header := []string{
		"Repository",
		"Author",
		"Closed",
		"ClosedAt",
		"CreatedAt",
		"LastEditedAt",
		"State",
		"Title",
		"UpdatedAt",
	}
	utils.WriteLineToFile(f, header...)

	for _, repo := range repos {
		for _, issue := range repo.Issues {
			cells := []string{
				repo.Name,
				issue.Author.Login,
				strconv.FormatBool(issue.Closed),
				issue.ClosedAt,
				issue.CreatedAt,
				issue.LastEditedAt,
				issue.State,
				strings.Replace(issue.Title, "\"", "'", -1),
				issue.UpdatedAt,
			}

			utils.WriteLineToFile(f, cells...)
		}
	}
}
