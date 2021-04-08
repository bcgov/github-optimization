package plugin

import (
	"context"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func (s *Server) handleOrganizationQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	query := &models.OrganizationQuery{}
	if err := UnmarshalQuery(q.JSON, query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(s.Datasource.HandleOrganizationQuery(ctx, query, q))
}

// HandleOrganization handles the plugin query for github Organization
func (s *Server) HandleOrganization(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{
		Responses: processQueries(ctx, req, s.handleOrganizationQuery),
	}, nil
}
