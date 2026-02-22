package ai

import "context"

type AIClient interface {
	ChatCompletion(ctx context.Context, prompt string) (string, error)
}