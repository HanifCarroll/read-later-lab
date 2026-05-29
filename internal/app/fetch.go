package app

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	titleRE       = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	metaRE        = regexp.MustCompile(`(?is)<meta\s+([^>]+)>`)
	attrRE        = regexp.MustCompile(`(?is)([a-zA-Z:-]+)\s*=\s*["']([^"']*)["']`)
	scriptStyleRE = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>|<style[^>]*>.*?</style>|<noscript[^>]*>.*?</noscript>`)
	tagRE         = regexp.MustCompile(`(?is)<[^>]+>`)
	spaceRE       = regexp.MustCompile(`\s+`)
)

func FetchArticle(ctx context.Context, rawURL string) (Article, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return Article{}, err
	}
	req.Header.Set("User-Agent", "read-later-experiment/0.1")

	res, err := client.Do(req)
	if err != nil {
		return Article{}, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Article{}, fmt.Errorf("fetch article: %s", res.Status)
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, 2_000_000))
	if err != nil {
		return Article{}, err
	}
	doc := string(body)
	metadata := extractMetadata(doc)
	u, _ := url.Parse(rawURL)

	article := Article{
		URL:         rawURL,
		Title:       first(metadata["og:title"], metadata["twitter:title"], extractTitle(doc)),
		Site:        first(metadata["og:site_name"], u.Hostname()),
		Author:      first(metadata["author"], metadata["article:author"]),
		PublishedAt: first(metadata["article:published_time"], metadata["date"], metadata["pubdate"]),
		Text:        extractText(doc),
	}
	if article.Text == "" {
		article.Text = article.Title
	}
	return article, nil
}

func extractMetadata(doc string) map[string]string {
	metadata := map[string]string{}
	for _, match := range metaRE.FindAllStringSubmatch(doc, -1) {
		attrs := map[string]string{}
		for _, attr := range attrRE.FindAllStringSubmatch(match[1], -1) {
			attrs[strings.ToLower(attr[1])] = html.UnescapeString(attr[2])
		}
		key := strings.ToLower(first(attrs["property"], attrs["name"], attrs["itemprop"]))
		if key != "" && attrs["content"] != "" {
			metadata[key] = attrs["content"]
		}
	}
	return metadata
}

func extractTitle(doc string) string {
	if match := titleRE.FindStringSubmatch(doc); len(match) == 2 {
		return cleanText(match[1])
	}
	return ""
}

func extractText(doc string) string {
	cleaned := scriptStyleRE.ReplaceAllString(doc, " ")
	cleaned = tagRE.ReplaceAllString(cleaned, " ")
	text := cleanText(cleaned)
	if len(text) > 30000 {
		return text[:30000]
	}
	return text
}

func cleanText(value string) string {
	value = html.UnescapeString(value)
	value = spaceRE.ReplaceAllString(value, " ")
	return strings.TrimSpace(value)
}

func first(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
