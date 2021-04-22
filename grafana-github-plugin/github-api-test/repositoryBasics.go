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

func writeLineToFile(f *os.File, cells [23]string) {
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
	Name  string
	Owner struct {
		Login string
	}
	NameWithOwner      string
	URL                string
	HomepageURL        string
	Description        string
	ForkCount          int
	IsFork             bool
	IsMirror           bool
	IsPrivate          bool
	IsArchived         bool
	IsTemplate         bool
	StargazerCount     int
	DiskUsage          int
	HasIssuesEnabled   bool
	HasProjectsEnabled bool
	HasWikiEnabled     bool
	MergeCommitAllowed bool
	RebaseMergeAllowed bool
	SquashMergeAllowed bool
	CreatedAt          string
	UpdatedAt          string
	PushedAt           string
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

		repos  = []Repository{}
		page   = 1
	)

	for {
		fmt.Printf("Querying %v page...\n", page)

		query := &QueryListRepositories{}
		if err := client.Query(ctx, query, variables); err != nil {
			fmt.Println(err)
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

		// time.Sleep(10 * time.Minute)
	}

	return repos, nil
}

func main() {
	// Try creating csv first
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	targetDir := "../../../notebook/dat/"
	targetFile := "/repository-basics.csv"

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

	repos, _ := GetRepositories(context.Background(), client)

	// Append data into csv
	header := [...]string{
		"Repository",
		"Owner",
		"Name With Owner",
		"Url",
		"Homepage Url",
		"Description",
		"Forks",
		"Is Fork",
		"Is Mirror",
		"Is Private",
		"Is Archived",
		"Is Template",
		"Stars",
		"Disk Usage",
		"Has Issues Enabled",
		"Has Projects Enabled",
		"Has Wiki Enabled",
		"Merge Commit Allowed",
		"Rebase Merge Allowed",
		"Squash Merge Allowed",
		"Created At",
		"Updated At",
		"Pushed At",
	}
	writeLineToFile(f, header)

	for _, repo := range repos {
		cells := [...]string{
			repo.Name,
			repo.Owner.Login,
			repo.NameWithOwner,
			repo.URL,
			repo.HomepageURL,
			repo.Description,
			strconv.Itoa(repo.StargazerCount),
			strconv.Itoa(repo.ForkCount),
			strconv.Itoa(repo.DiskUsage),
			strconv.FormatBool(repo.IsFork),
			strconv.FormatBool(repo.IsMirror),
			strconv.FormatBool(repo.IsPrivate),
			strconv.FormatBool(repo.IsArchived),
			strconv.FormatBool(repo.IsTemplate),
			strconv.FormatBool(repo.HasIssuesEnabled),
			strconv.FormatBool(repo.HasProjectsEnabled),
			strconv.FormatBool(repo.HasWikiEnabled),
			strconv.FormatBool(repo.MergeCommitAllowed),
			strconv.FormatBool(repo.RebaseMergeAllowed),
			strconv.FormatBool(repo.SquashMergeAllowed),
			repo.CreatedAt,
			repo.UpdatedAt,
			repo.PushedAt,
		}

		writeLineToFile(f, cells)
	}
}
