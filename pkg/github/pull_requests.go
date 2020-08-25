package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryListPullRequests lists all pull requests in a repository
// {
//   search(query: "is:pr repo:grafana/grafana merged:2020-08-19..*", type: ISSUE, first: 100) {
//     nodes {
//       ... on PullRequest {
//         id
//         title
//       }
//   }
// }
type QueryListPullRequests struct {
	Search struct {
		Nodes []struct {
			PullRequest PullRequest `graphql:"... on PullRequest"`
		}
		PageInfo PageInfo
	} `graphql:"search(query: $query, type: ISSUE, first: 100, after: $cursor)"`
}

type PullRequestAuthor struct {
	User User `graphql:"... on User"`
}

// PullRequest is a GitHub pull request
type PullRequest struct {
	Title     string
	State     githubv4.PullRequestState
	Author    PullRequestAuthor
	Closed    bool
	IsDraft   bool
	Locked    bool
	Merged    bool
	ClosedAt  githubv4.DateTime
	CreatedAt githubv4.DateTime
	UpdatedAt githubv4.DateTime
	MergedAt  githubv4.DateTime
	Mergeable githubv4.MergeableState
	MergedBy  *PullRequestAuthor
}

// PullRequests is a list of GitHub Pull Requests
type PullRequests []PullRequest

// Frame converts the list of Pull Requests to a Grafana DataFrame
func (p PullRequests) Frame() data.Frames {
	frame := data.NewFrame(
		"pull_requests",
		data.NewField("title", nil, []string{}),
		data.NewField("state", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("author_email", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("closed", nil, []bool{}),
		data.NewField("is_draft", nil, []bool{}),
		data.NewField("locked", nil, []bool{}),
		data.NewField("merged", nil, []bool{}),
		data.NewField("mergeable", nil, []string{}),
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("merged_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("created_at", nil, []time.Time{}),
	)

	for _, v := range p {
		var (
			closedAt *time.Time
			mergedAt *time.Time
		)

		if !v.ClosedAt.IsZero() {
			closedAt = &v.ClosedAt.Time
		}

		if !v.MergedAt.IsZero() {
			mergedAt = &v.MergedAt.Time
		}

		frame.AppendRow(
			v.Title,
			string(v.State),
			v.Author.User.Login,
			v.Author.User.Email,
			v.Author.User.Company,
			v.Closed,
			v.IsDraft,
			v.Locked,
			v.Merged,
			string(v.Mergeable),
			closedAt,
			mergedAt,
			v.UpdatedAt.Time,
			v.CreatedAt.Time,
		)
	}

	return data.Frames{frame}

}

// GetPullRequestsInRange uses the graphql search endpoint API to find pull requests in the given time range.
func GetPullRequestsInRange(ctx context.Context, client Client, opts models.ListPullRequestsInRangeOptions, from time.Time, to time.Time) (PullRequests, error) {
	search := []string{
		"is:pr",
		fmt.Sprintf("repo:%s/%s", opts.Owner, opts.Repository),
		fmt.Sprintf("%s:%s..%s", opts.TimeField.String(), from.Format(time.RFC3339), to.Format(time.RFC3339)),
	}

	if opts.Query != nil {
		search = append(search, *opts.Query)
	}

	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(strings.Join(search, " ")),
		}

		pullRequests = []PullRequest{}
	)

	for {
		q := &QueryListPullRequests{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}
		prs := make([]PullRequest, len(q.Search.Nodes))

		for i, v := range q.Search.Nodes {
			prs[i] = v.PullRequest
		}

		pullRequests = append(pullRequests, prs...)

		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Search.PageInfo.EndCursor
	}

	return pullRequests, nil
}
