package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Store interface {
	List(context.Context) ([]Item, error)
	Create(context.Context, Item) (Item, error)
	UpdateStatus(context.Context, string, ItemStatus) (Item, error)
}

type Analyzer interface {
	Analyze(context.Context, Article, string) (Analysis, error)
}

type Service struct {
	store    Store
	analyzer Analyzer
}

func NewService(store Store, analyzer Analyzer) *Service {
	return &Service{store: store, analyzer: analyzer}
}

type CreateItemInput struct {
	URL         string `json:"url"`
	SavedReason string `json:"savedReason"`
}

func (s *Service) ListItems(ctx context.Context) ([]Item, error) {
	items, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (s *Service) CreateItem(ctx context.Context, input CreateItemInput) (Item, error) {
	input.URL = strings.TrimSpace(input.URL)
	input.SavedReason = strings.TrimSpace(input.SavedReason)
	if input.URL == "" {
		return Item{}, errors.New("url is required")
	}
	if _, err := url.ParseRequestURI(input.URL); err != nil {
		return Item{}, fmt.Errorf("valid url is required")
	}
	if input.SavedReason == "" {
		return Item{}, errors.New("saved reason is required")
	}

	article, err := FetchArticle(ctx, input.URL)
	if err != nil {
		return Item{}, err
	}
	analysis, err := s.analyzer.Analyze(ctx, article, input.SavedReason)
	if err != nil {
		return Item{}, err
	}

	now := time.Now().UTC()
	item := Item{
		ID:              newID(),
		URL:             article.URL,
		Title:           fallback(article.Title, input.URL),
		Site:            article.Site,
		Author:          article.Author,
		PublishedAt:     article.PublishedAt,
		SavedReason:     input.SavedReason,
		Summary:         analysis.Summary,
		ReadTimeMinutes: analysis.ReadTimeMinutes,
		RecommendedMode: normalizeMode(analysis.RecommendedMode),
		BestSections:    trimTo(analysis.BestSections, 4),
		ReadIf:          trimTo(analysis.ReadIf, 3),
		SkipIf:          trimTo(analysis.SkipIf, 3),
		Status:          StatusInbox,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	return s.store.Create(ctx, item)
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status ItemStatus) (Item, error) {
	switch status {
	case StatusInbox, StatusReadSoon, StatusSkimLater, StatusReference, StatusSkipped, StatusArchived:
		return s.store.UpdateStatus(ctx, id, status)
	default:
		return Item{}, errors.New("unknown status")
	}
}

func newID() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}

func fallback(v, fb string) string {
	if strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v)
	}
	return fb
}

func trimTo(values []string, n int) []string {
	out := make([]string, 0, n)
	for _, v := range values {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
		if len(out) == n {
			break
		}
	}
	return out
}

func normalizeMode(mode ReadingMode) ReadingMode {
	switch mode {
	case ModeReadFully, ModeSkim, ModeReference, ModeSkip:
		return mode
	default:
		return ModeSkim
	}
}
