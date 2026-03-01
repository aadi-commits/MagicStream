package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

func buildAdminRankingPrompt(review string) string {
	return fmt.Sprintf(`You are a strict movie review ranking classifier.

		Analyze the following admin-written movie review and assign exactly ONE ranking number based on sentiment.

		Ranking scale:
		1 = Excellent (very positive, outstanding, exceptional)
		2 = Good (positive, enjoyable, well done)
		3 = Okay (average, neutral, mixed feelings)
		4 = Bad (negative, disappointing, poor quality)
		5 = Terrible (very negative, awful, extremely poor)

		STRICT RULES:
		- Return ONLY ONE number.
		- Allowed outputs: 1 or 2 or 3 or 4 or 5.
		- Do NOT return text.
		- Do NOT return explanation.
		- Do NOT return JSON.
		- Do NOT return words.
		- Do NOT return multiple numbers.
		- Do NOT include punctuation.
		- Output must contain only the digit.

		Admin Review:
		"""
		%s
		"""

		Now return ONLY the number.`, review)
}

func (s *Service) GenerateAdminRanking(ctx context.Context, review string) (int, string, error) {

	prompt := buildAdminRankingPrompt(review)

	rawResponse, err := s.client.ChatCompletion(ctx, prompt)
	if err != nil {
		return 0, rawResponse, fmt.Errorf("AI request failed: %w", err)
	}

	fmt.Println("RAW AI RESPONSE:", rawResponse)
	rankStr := strings.TrimSpace(rawResponse)

	switch rankStr {
	case "1", "2", "3", "4", "5":
		return int(rankStr[0] - '0'), rawResponse, nil
	default:
		return 0, rawResponse, errors.New("Invalid AI ranking response")
	}
}