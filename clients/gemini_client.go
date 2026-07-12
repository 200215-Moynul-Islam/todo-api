package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

const (
	defaultGeminiModel = "gemini-flash-latest"
	geminiEndpoint     = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent"
)

// GeminiClient communicates with the Gemini API.
type GeminiClient interface {
	GenerateDescription(title string) (string, error)
}

type geminiClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

// NewGeminiClient creates a client using Gemini settings from app.conf.
func NewGeminiClient() GeminiClient {
	apiKey, _ := beego.AppConfig.String("GEMINI_API_KEY")
	model, _ := beego.AppConfig.String("GEMINI_MODEL")
	if strings.TrimSpace(model) == "" {
		model = defaultGeminiModel
	}

	return &geminiClient{
		apiKey: strings.TrimSpace(apiKey),
		model:  strings.TrimSpace(model),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *geminiClient) GenerateDescription(title string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	payload := geminiRequest{
		Contents: []geminiContent{{
			Parts: []geminiPart{{
				Text: "Generate a helpful task description related to the title in no more than two sentences. Return only the description. Task title: " + title,
			}},
		}},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("create Gemini request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(geminiEndpoint, c.model), bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create Gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("Gemini API returned status %d", resp.StatusCode)
	}

	var result geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode Gemini response: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini API returned no description")
	}

	description := limitToTwoSentences(result.Candidates[0].Content.Parts[0].Text)
	if description == "" {
		return "", fmt.Errorf("Gemini API returned an empty description")
	}
	return description, nil
}

func limitToTwoSentences(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	end := 0
	for i, char := range text {
		if char == '.' || char == '!' || char == '?' {
			end++
			if end == 2 {
				return strings.TrimSpace(text[:i+1])
			}
		}
	}
	return text
}
