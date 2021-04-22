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

func writeLineToFile(f *os.File, cells [9]string) {
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
	Packages struct {
		TotalCount int
	}
	Projects struct {
		TotalCount int
	}
	Releases struct {
		TotalCount int
	}
	Submodules struct {
		TotalCount int
	}
	DeployKeys struct {
		TotalCount int
	}
	RepositoryTopics struct {
		TotalCount int
	}
	LicenseInfo struct {
		Name string
	}
	CodeOfConduct struct {
		Name string
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
	targetFile := "/repository-details-1.csv"

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
		"PackageCount",
		"ProjectCount",
		"ReleaseCount",
		"SubmoduleCount",
		"DeployKeyCount",
		"TopicCount",
		"License",
		"CodeOfConduct",
	}
	writeLineToFile(f, header)

	for _, repo := range repos {
		name := repo.Name
		packageCount := repo.Packages.TotalCount
		projectCount := repo.Projects.TotalCount
		releaseCount := repo.Releases.TotalCount
		submoduleCount := repo.Submodules.TotalCount
		deployKeyCount := repo.DeployKeys.TotalCount
		topicCount := repo.RepositoryTopics.TotalCount
		license := repo.LicenseInfo.Name
		codeOfConduct := repo.CodeOfConduct.Name

		packageCountStr := strconv.Itoa(packageCount)
		projectCountStr := strconv.Itoa(projectCount)
		releaseCountStr := strconv.Itoa(releaseCount)
		submoduleCountStr := strconv.Itoa(submoduleCount)
		deployKeyCountStr := strconv.Itoa(deployKeyCount)
		topicCountStr := strconv.Itoa(topicCount)

		cells := [...]string{
			name,
			packageCountStr,
			projectCountStr,
			releaseCountStr,
			submoduleCountStr,
			deployKeyCountStr,
			topicCountStr,
			license,
			codeOfConduct,
		}

		writeLineToFile(f, cells)
	}
}
