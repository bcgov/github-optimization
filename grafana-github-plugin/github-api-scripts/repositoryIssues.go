package main

// Template for using the library.

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

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

type QueryListRepositoryIssues struct {
	Organization struct {
		Repositories struct {
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
			Nodes []struct {
				Name   githubv4.String
				Issues struct {
					Nodes []struct {
						Author struct {
							Login githubv4.String
						}
						Closed       githubv4.Boolean
						ClosedAt     githubv4.String
						CreatedAt    githubv4.String
						LastEditedAt githubv4.String
						State        githubv4.String
						Title        githubv4.String
						UpdatedAt    githubv4.String
					}
					PageInfo struct {
						HasNextPage githubv4.Boolean
						EndCursor   githubv4.String
					}
				} `graphql:"issues(first: 100)"`
			}
		} `graphql:"repositories(first: 100, after: $reposCursor)"`
	} `graphql:"organization(login: $org)"`
}

type QueryListPaginatedRepositoryIssues struct {
	Organization struct {
		Repository struct {
			Issues struct {
				Nodes []struct {
					Author struct {
						Login githubv4.String
					}
					Closed       githubv4.Boolean
					ClosedAt     githubv4.String
					CreatedAt    githubv4.String
					LastEditedAt githubv4.String
					State        githubv4.String
					Title        githubv4.String
					UpdatedAt    githubv4.String
				}
				PageInfo struct {
					HasNextPage githubv4.Boolean
					EndCursor   githubv4.String
				}
			} `graphql:"issues(first: 100, after: $issuesCursor)"`
		} `graphql:"repository(name: $name)"`
	} `graphql:"organization(login: $org)"`
}

func listPaginatedIssues(client *githubv4.Client, f *os.File, name githubv4.String, issuesCursor githubv4.String) {
	variables := map[string]interface{}{
		"org":          githubv4.String("bcgov"),
		"name":         name,
		"issuesCursor": issuesCursor,
	}
	for {
		q := &QueryListPaginatedRepositoryIssues{}
		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			fmt.Print(err)
		}
		for _, b := range q.Organization.Repository.Issues.Nodes {
			author := string(b.Author.Login)
			closed := strconv.FormatBool(bool(b.Closed))
			closedAt := string(b.ClosedAt)
			createdAt := string(b.CreatedAt)
			lastEditedAt := string(b.LastEditedAt)
			state := string(b.State)
			title := strings.Replace(string(b.Title), "\"", "'", -1)
			updatedAt := string(b.UpdatedAt)
			cells := [...]string{string(name), author, closed, closedAt, createdAt, lastEditedAt, state, "\"" + title + "\"", updatedAt}
			writeLineToFile(f, cells)
		}
		fmt.Println("looping issues for repo: " + string(name))
		if !q.Organization.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		variables["issuesCursor"] = q.Organization.Repository.Issues.PageInfo.EndCursor
	}
}

func main() {
	variables := map[string]interface{}{
		"org":         githubv4.String("bcgov"),
		"reposCursor": (*githubv4.String)(nil),
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

	f, err := os.Create(path + "/dat/issue-topics.csv")
	check(err)
	defer f.Close()

	f.WriteString("Repository,Author,Closed,ClosedAt, CreatedAt, LastEditedAt, State, Title, UpdatedAt\n")

	for {
		q := &QueryListRepositoryIssues{}
		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			fmt.Print(err)
		}
		for _, b := range q.Organization.Repositories.Nodes {
			name := b.Name
			if b.Issues.PageInfo.HasNextPage {
				listPaginatedIssues(client, f, name, b.Issues.PageInfo.EndCursor)
			} else {
				variables["reposCursor"] = q.Organization.Repositories.PageInfo.EndCursor
			}
			for _, c := range b.Issues.Nodes {
				author := string(c.Author.Login)
				closed := strconv.FormatBool(bool(c.Closed))
				closedAt := string(c.ClosedAt)
				createdAt := string(c.CreatedAt)
				lastEditedAt := string(c.LastEditedAt)
				state := string(c.State)
				title := strings.Replace(string(c.Title), "\"", "'", -1)
				updatedAt := string(c.UpdatedAt)
				cells := [...]string{string(name), author, closed, closedAt, createdAt, lastEditedAt, state, "\"" + title + "\"", updatedAt}
				writeLineToFile(f, cells)
			}
		}

		if !q.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
	}
}
