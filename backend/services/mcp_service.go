package services

import (
	"context"
	"time"
)

type MCPService struct {
	timeout time.Duration
}

func NewMCPService(timeout time.Duration) *MCPService {
	return &MCPService{
		timeout: timeout,
	}
}

func (s *MCPService) DoSomething(ctx context.Context) error {
	// 使用 context
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.timeout):
		return nil
	}
}

