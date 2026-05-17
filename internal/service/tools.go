package service

import (
	"context"
	"sync"

	v1 "tool-go/api/v1"
)

type ITools interface {
	MockData(ctx context.Context, req *v1.MockDataReq) (*v1.MockDataRes, error)
}

var (
	localTools ITools
	toolsMu    sync.RWMutex
)

func Tools() ITools {
	toolsMu.RLock()
	defer toolsMu.RUnlock()
	return localTools
}

func RegisterTools(i ITools) {
	toolsMu.Lock()
	defer toolsMu.Unlock()
	localTools = i
}
