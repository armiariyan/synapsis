package healthcheck

import "context"

type Service interface {
	HealthCheck(ctx context.Context) (res HealthCheckResponse, err error)
}