package postgres

import (
	"context"
	"fmt"

	"irj/pkg/dto"

	"github.com/jackc/pgx/v5/pgxpool"

	_http "irj/pkg/http"
)

type (
	Service struct {
		pool *pgxpool.Pool
	}
)

func (s *Service) Close() {
	s.pool.Close()
}

func (s *Service) Status(ctx context.Context) *dto.ServiceItem {
	status := _http.OK
	details := ""

	var version string
	if err := s.pool.QueryRow(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		status = _http.KO
		details = fmt.Sprintf("failed to query version: %s", err)
	}

	return &dto.ServiceItem{
		Name:    "db",
		Status:  status,
		Details: details,
	}
}
