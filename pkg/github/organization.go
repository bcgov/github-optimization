package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

type Org struct {
	ID           githubv4.ID
	Repositories struct {
		TotalCount int64
	} `graphql:"repositories"`
	Packages struct {
		TotalCount int64
	} `graphql:"packages"`
	Projects struct {
		TotalCount int64
	} `graphql:"projects"`
}

type QueryOrganization struct {
	Viewer struct {
		Organization Org `graphql:"organization(login: $org)"`
	}
}

// Frames converts the Organization to a Grafana DataFrame
func (c Org) Frames() data.Frames {
	frame := data.NewFrame(
		"organization",
		data.NewField("id", nil, []string{}),
		data.NewField("repo_count", nil, []int64{}),
		data.NewField("package_count", nil, []int64{}),
		data.NewField("project_count", nil, []int64{}),
	)

	frame.AppendRow(
		c.ID,
		c.Repositories.TotalCount,
		c.Packages.TotalCount,
		c.Projects.TotalCount,
	)

	return data.Frames{frame}
}

// GetOrganization retruns the organization basic information for the client
func GetOrganization(ctx context.Context, client Client, opts models.GetOrganizationOptions) (Org, error) {
	query := &QueryOrganization{}

	variables := map[string]interface{}{
		"org": githubv4.String(opts.Owner),
	}

	if err := client.Query(ctx, query, variables); err != nil {
		return Org{}, err
	}

	return query.Viewer.Organization, nil
}
