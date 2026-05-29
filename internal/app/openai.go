package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type OpenAIAnalyzer struct {
	APIKey string
	Model  string
}

func (a OpenAIAnalyzer) Analyze(ctx context.Context, article Article, savedReason string) (Analysis, error) {
	if strings.TrimSpace(a.APIKey) == "" {
		return fallbackAnalysis(article), nil
	}

	payload := map[string]any{
		"model":           a.Model,
		"response_format": map[string]string{"type": "json_object"},
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You analyze saved articles for a read-it-later inbox. Return compact JSON only with keys summary, readTimeMinutes, recommendedMode, bestSections, readIf, skipIf. recommendedMode must be one of read_fully, skim, reference, skip. bestSections is 2-4 concise strings. readIf and skipIf are exactly 3 concise bullet strings each.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Saved reason: %s\n\nTitle: %s\nSite: %s\nURL: %s\n\nArticle text:\n%s", savedReason, article.Title, article.Site, article.URL, article.Text),
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return Analysis{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Analysis{}, err
	}
	req.Header.Set("Authorization", "Bearer "+a.APIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Analysis{}, err
	}
	defer res.Body.Close()

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return Analysis{}, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return Analysis{}, errors.New(decoded.Error.Message)
		}
		return Analysis{}, fmt.Errorf("openai request failed: %s", res.Status)
	}
	if len(decoded.Choices) == 0 {
		return Analysis{}, errors.New("openai returned no choices")
	}

	var analysis Analysis
	if err := json.Unmarshal([]byte(decoded.Choices[0].Message.Content), &analysis); err != nil {
		return Analysis{}, err
	}
	analysis.RecommendedMode = normalizeMode(analysis.RecommendedMode)
	if analysis.ReadTimeMinutes <= 0 {
		analysis.ReadTimeMinutes = estimateReadTime(article.Text)
	}
	return analysis, nil
}

func fallbackAnalysis(article Article) Analysis {
	return Analysis{
		Summary:         "LLM analysis is unavailable because OPENAI_API_KEY is not set; this placeholder uses fetched article metadata.",
		ReadTimeMinutes: estimateReadTime(article.Text),
		RecommendedMode: ModeSkim,
		BestSections:    []string{"Review the introduction", "Scan headings for relevant sections"},
		ReadIf:          []string{"The title matches a current interest", "You want a quick overview", "You are collecting references on this topic"},
		SkipIf:          []string{"You need verified deep analysis", "The source is not relevant", "You only saved it out of habit"},
	}
}

func estimateReadTime(text string) int {
	words := len(strings.Fields(text))
	minutes := words / 225
	if words%225 != 0 {
		minutes++
	}
	if minutes < 1 {
		return 1
	}
	return minutes
}
