package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type QueryListRepositoryTopics struct {
	Organization struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
			Nodes []struct {
				Name             githubv4.String
				RepositoryTopics struct {
					Nodes []struct {
						Topic struct {
							Name githubv4.String
						}
					}
				} `graphql:"repositoryTopics(first: 100)"`
			}
		} `graphql:"repositories(first: 100, after: $cursor)"`
	} `graphql:"organization(login: $org)"`
}

func main() {
	codes := [...]string{"AEST", "AGRI", "ALC", "AG", "MCF", "CITZ", "DBC", "EMBC", "EAO", "EDUC", "EMPR", "ENV", "FIN", "FLNR", "HLTH", "IRR", "JEDC", "LBR", "LDB", "MMHA", "MAH", "BCPC", "PSA", "PSSG", "SDPR", "TCA", "TRAN"}

	variables := map[string]interface{}{
		"org":    githubv4.String("bcgov"),
		"cursor": (*githubv4.String)(nil),
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(path + "/dat/repo-topics.txt")
	check(err)
	defer f.Close()

	for {
		q := &QueryListRepositoryTopics{}
		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			fmt.Print(err)
		}
		for _, b := range q.Organization.Repositories.Nodes {
			name := b.Name
			f.WriteString(string(name) + ": ")
			// topics := make([]string, len(b.RepositoryTopics.Nodes))
			for _, c := range b.RepositoryTopics.Nodes {
				if contains(codes[:], strings.ToUpper(string(c.Topic.Name))) {
					f.WriteString(string(c.Topic.Name) + ", ")
					// topics = append(topics, string(c.Topic.Name))
				}
			}
			f.WriteString("\n")

		}
		if !q.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Organization.Repositories.PageInfo.EndCursor
	}

}
