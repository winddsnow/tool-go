package controller

import (
	"context"

	v1 "tool-go/api/v1"
	"tool-go/internal/service"
)

var Tools = cTools{}

type cTools struct{}

func (c *cTools) MockData(ctx context.Context, req *v1.MockDataReq) (*v1.MockDataRes, error) {
	return service.Tools().MockData(ctx, req)
}
