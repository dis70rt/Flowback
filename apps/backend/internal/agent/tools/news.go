package tools

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

type SearchNewsInput struct {
	Location string `json:"location"`
	Query    string `json:"query"`
}

type SearchNewsOutput struct {
	Headlines []string `json:"headlines"`
	Summary   string   `json:"summary"`
}

func FetchLocalNews(location string, query string) ([]string, string, error) {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		slog.Warn("news context skipped: NEWS_API_KEY not set", "city", location)
		return nil, "", fmt.Errorf("NEWS_API_KEY is not configured in the environment")
	}

	slog.Info("fetching local news context", "city", location, "query", query)

	searchQuery := fmt.Sprintf("%s %s", location, query)
	apiURL := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&sortBy=publishedAt&apiKey=%s",
		url.QueryEscape(searchQuery), apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to contact news api: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("news api returned status: %d", resp.StatusCode)
	}

	var result struct {
		Status   string `json:"status"`
		Articles []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"articles"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", fmt.Errorf("failed to parse news api response: %v", err)
	}

	var headlines []string
	for i, article := range result.Articles {
		if i >= 3 {
			break
		}
		headlines = append(headlines, article.Title)
	}

	summary := "No major disruptions found."
	if len(headlines) > 0 {
		slog.Info("news context injected into payload", "city", location, "headline_count", len(headlines), "headlines", headlines)
		summary = "Recent news articles suggest potential disruptions related to the query."
	} else {
		slog.Info("no relevant news headlines found for city", "city", location)
	}

	return headlines, summary, nil
}

func NewSearchLocalNewsTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "search_local_news",
			Description: "Searches the web for breaking news in the customer's location (e.g., floods, banking outages) that might explain a payment failure.",
		},
		func(ctx agent.Context, input SearchNewsInput) (*SearchNewsOutput, error) {

			headlines, summary, err := FetchLocalNews(input.Location, input.Query)
			if err != nil {
				return nil, err
			}
			return &SearchNewsOutput{
				Headlines: headlines,
				Summary:   summary,
			}, nil
		},
	)
}
