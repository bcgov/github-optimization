package main

// Template for using the library.

import (
	"context"
	"fmt"
	"math"
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

func writeLineToFile(f *os.File, cells [10]string) {
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
	CreatedAt githubv4.DateTime
	UpdatedAt githubv4.DateTime
	PushedAt  githubv4.DateTime
	Issues    struct {
		TotalCount int
	}
	PullRequests struct {
		TotalCount int
	}
	DefaultBranchRef struct {
		Name   string
		Prefix string
	}
	Languages struct {
		Nodes []struct {
			Name string
		}
	} `graphql:"languages(first: 100, orderBy: {field: SIZE, direction: DESC})"`
}

// Repositories is a list of GitHub repositories
type Repositories []Repository

// RepositoryExtra is ...
type RepositoryExtra struct {
	DefaultBranchCommitCount int
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
				RefName: v.DefaultBranchRef.Name,
			}

			n, _ := GetDefaultBranchCommitCount(ctx, client, opts)
			extra := RepositoryExtra{
				DefaultBranchCommitCount: n,
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
	RefName string
}

// QueryDefaultBranchCommitCount is the GraphQL query for retrieving a list of repositories for an organization
// query {
// 	 repository(name: "$name", owner: "$owner") {
//     ref(qualifiedName: "$qualifiedName") {
// 	     target {
// 	       ... on Commit {
//           history {
//             totalCount
//           }
//         }
//       }
// 	   }
// 	 }
// }
type QueryDefaultBranchCommitCount struct {
	Repository struct {
		Ref struct {
			Target struct {
				Commit struct {
					History struct {
						TotalCount int
					}
				} `graphql:"... on Commit"`
			}
		} `graphql:"ref(qualifiedName: $refName)"`
	} `graphql:"repository(name: $name, owner: $org)"`
}

// GetDefaultBranchCommitCount retruns ...
func GetDefaultBranchCommitCount(ctx context.Context, client Client, ops BranchOptions) (int, error) {
	var (
		variables = map[string]interface{}{
			"org":     githubv4.String(ops.Org),
			"name":    githubv4.String(ops.Name),
			"refName": githubv4.String(ops.RefName),
		}
	)

	query := &QueryDefaultBranchCommitCount{}
	if err := client.Query(ctx, query, variables); err != nil {
		fmt.Println(err)
		return 0, err
	}

	count := query.Repository.Ref.Target.Commit.History.TotalCount

	return count, nil
}

func main() {
	// Try creating csv first
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	targetDir := "../../../notebook/dat/"
	targetFile := "/repository-details-2.csv"

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
		"DaysOpen",
		"IssueCount",
		"PrCount",
		"CommitCount",
		"AverageIssueCountPerDay",
		"AveragePrCountPerDay",
		"AverageCommitCountPerDay",
		"DefaultBranchName",
		"Languages",
	}
	writeLineToFile(f, header)

	for i, repo := range repos {
		name := repo.Name
		issueCount := repo.Issues.TotalCount
		prCount := repo.PullRequests.TotalCount
		commitCount := extras[i].DefaultBranchCommitCount
		defaultBranchName := repo.DefaultBranchRef.Name

		languages := ""
		for i, v := range repo.Languages.Nodes {
			if i > 0 {
				languages += "_"
			}
			languages += v.Name
		}

		hoursOpen := time.Now().UTC().Sub(repo.CreatedAt.UTC()).Hours()
		daysOpen := hoursOpen / 24

		averageIssueCountPerDay := float64(issueCount) / daysOpen
		averagePrCountPerDay := float64(prCount) / daysOpen
		averageCommitCountPerDay := float64(commitCount) / daysOpen

		issueCountStr := strconv.Itoa(issueCount)
		prCountStr := strconv.Itoa(prCount)
		commitCountStr := strconv.Itoa(commitCount)
		daysOpenStr := strconv.Itoa(int(daysOpen))
		averageIssueCountPerDayStr := strconv.FormatFloat(math.Round(averageIssueCountPerDay*100)/100, 'f', -1, 32)
		averagePrCountPerDayStr := strconv.FormatFloat(math.Round(averagePrCountPerDay*100)/100, 'f', -1, 32)
		averageCommitCountPerDayStr := strconv.FormatFloat(math.Round(averageCommitCountPerDay*100)/100, 'f', -1, 32)

		cells := [...]string{
			name,
			daysOpenStr,
			issueCountStr,
			prCountStr,
			commitCountStr,
			averageIssueCountPerDayStr,
			averagePrCountPerDayStr,
			averageCommitCountPerDayStr,
			defaultBranchName,
			languages,
		}

		writeLineToFile(f, cells)
	}
}
