package ai

import (
	"context"
)

type Service struct {
	client AIClient
}

func NewService(client AIClient) *Service{
	return &Service{
		client: client,
	}
}

func (s *Service) GenerateReply(ctx context.Context, prompt string) (string, error) {

	return s.client.ChatCompletion(ctx, prompt)
}