package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

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
	Search struct {
		Nodes []struct {
			Repository Repository `graphql:"... on Repository"`
		}
		PageInfo PageInfo
	} `graphql:"search(query: $query, type: REPOSITORY, first: 100, after: $cursor)"`
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
	ForkCount          int64
	IsFork             bool
	IsMirror           bool
	IsPrivate          bool
	IsArchived         bool
	IsTemplate         bool
	StargazerCount     int64
	DiskUsage          int64
	HasIssuesEnabled   bool
	HasProjectsEnabled bool
	HasWikiEnabled     bool
	MergeCommitAllowed bool
	RebaseMergeAllowed bool
	SquashMergeAllowed bool
	CreatedAt          githubv4.DateTime
	UpdatedAt          githubv4.DateTime
	PushedAt           githubv4.DateTime
}

// Repositories is a list of GitHub repositories
type Repositories []Repository

// Frames converts the list of GitHub repositories to a Grafana Dataframe
func (r Repositories) Frames() data.Frames {
	frame := data.NewFrame(
		"repositories",
		data.NewField("name", nil, []string{}),
		data.NewField("owner", nil, []string{}),
		data.NewField("name_with_owner", nil, []string{}),
		data.NewField("url", nil, []string{}),
		data.NewField("homepage_url", nil, []string{}),
		data.NewField("description", nil, []string{}),
		data.NewField("forks", nil, []int64{}),
		data.NewField("is_fork", nil, []bool{}),
		data.NewField("is_mirror", nil, []bool{}),
		data.NewField("is_private", nil, []bool{}),
		data.NewField("is_archived", nil, []bool{}),
		data.NewField("is_template", nil, []bool{}),
		data.NewField("stars", nil, []int64{}),
		data.NewField("disk_usage", nil, []int64{}),
		data.NewField("has_issues_enabled", nil, []bool{}),
		data.NewField("has_projects_enabled", nil, []bool{}),
		data.NewField("has_wiki_enabled", nil, []bool{}),
		data.NewField("merge_commit_allowed", nil, []bool{}),
		data.NewField("rebase_merge_allowed", nil, []bool{}),
		data.NewField("squash_merge_allowed", nil, []bool{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("pushed_at", nil, []time.Time{}),
	)

	for _, v := range r {
		frame.AppendRow(
			v.Name,
			v.Owner.Login,
			v.NameWithOwner,
			v.URL,
			v.HomepageURL,
			v.Description,
			v.ForkCount,
			v.IsFork,
			v.IsMirror,
			v.IsPrivate,
			v.IsArchived,
			v.IsTemplate,
			v.StargazerCount,
			v.DiskUsage,
			v.HasIssuesEnabled,
			v.HasProjectsEnabled,
			v.HasWikiEnabled,
			v.MergeCommitAllowed,
			v.RebaseMergeAllowed,
			v.SquashMergeAllowed,
			v.CreatedAt.Time,
			v.UpdatedAt.Time,
			v.PushedAt.Time,
		)
	}

	return data.Frames{frame}

}

// GetAllRepositories retrieves all available repositories for an organization
func GetAllRepositories(ctx context.Context, client Client, opts models.ListRepositoriesOptions) (Repositories, error) {
	query := strings.Join([]string{
		fmt.Sprintf("org:%s", opts.Owner),
		opts.Repository,
	}, " ")

	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(query),
		}

		repos = []Repository{}
	)

	for {
		q := &QueryListRepositories{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		r := make([]Repository, len(q.Search.Nodes))

		for i, v := range q.Search.Nodes {
			r[i] = v.Repository
		}

		repos = append(repos, r...)

		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Search.PageInfo.EndCursor
	}

	return repos, nil
}
