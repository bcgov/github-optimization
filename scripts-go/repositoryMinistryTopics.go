package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"
	"strings"

	utils "gh.com/api-test/utils"
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

// Repository is a code repository
type Repository struct {
	Name             string
	RepositoryTopics struct {
		Nodes []struct {
			Topic struct {
				Name string
			}
		}
	} `graphql:"repositoryTopics(first: 100)"`
}

// Repositories is a list of GitHub repositories
type Repositories []Repository

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

		repos = []Repository{}
		page  = 1
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

var ministryCodes = [...]string{"AEST", "AGRI", "ALC", "AG", "MCF", "CITZ", "DBC", "EMBC", "EAO", "EDUC", "EMPR", "ENV", "FIN", "FLNR", "HLTH", "IRR", "JEDC", "LBR", "LDB", "MMHA", "MAH", "BCPC", "PSA", "PSSG", "SDPR", "TCA", "TRAN"}

func main() {
	token, org := utils.CheckEnv()

	// Try creating csv first
	path, err := os.Getwd()
	utils.HandleError(err)

	targetDir := fmt.Sprintf("../../notebook/dat/%v/", org)
	targetFile := "/repository-ministry-topics.csv"

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
		"repository",
		"ministry_code",
	}
	utils.WriteLineToFile(f, header...)

	for _, repo := range repos {
		hasTopic := false
		name := repo.Name
		for _, node := range repo.RepositoryTopics.Nodes {
			topic := strings.ToUpper(node.Topic.Name)

			if contains(ministryCodes[:], topic) {
				utils.WriteLineToFile(f, name, topic)
				hasTopic = true
			}
		}
		if !hasTopic {
			utils.WriteLineToFile(f, name, "")
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
