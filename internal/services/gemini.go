package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go-chat/internal/config"
)

// GeminiService handles communication with the Gemini API
type GeminiService struct {
	config *config.Config
	client *http.Client
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// LLMResponse represents the structured response from the LLM
type LLMResponse struct {
	Text    string `json:"text,omitempty"`
	Table   string `json:"table,omitempty"`
	Diagram string `json:"diagram,omitempty"`
}

func NewGeminiService(config *config.Config) *GeminiService {
	return &GeminiService{
		config: config,
		client: &http.Client{},
	}
}

func (s *GeminiService) GenerateResponse(prompt string) (LLMResponse, error) {
	// Prepare URL
	u, err := s.prepareURL()
	if err != nil {
		return LLMResponse{}, err
	}
	prompt = "Keep the response to max of 5 lines" + prompt
	// Send request
	resp, err := s.sendRequest(u, prompt)
	if err != nil {
		return LLMResponse{}, err
	}
	defer resp.Body.Close()

	// Handle response
	return s.processResponse(resp)
}

func (s *GeminiService) prepareURL() (*url.URL, error) {
	u, err := url.Parse(s.config.GeminiAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("key", s.config.GeminiAPIKey)
	u.RawQuery = q.Encode()
	return u, nil
}

func (s *GeminiService) sendRequest(u *url.URL, prompt string) (*http.Response, error) {
	reqBody := GeminiRequest{
		Contents: []Content{{
			Parts: []Part{{Text: prompt}},
		}},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return s.client.Do(req)
}

func (s *GeminiService) processResponse(resp *http.Response) (LLMResponse, error) {
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return LLMResponse{}, fmt.Errorf("API error (status %d)", resp.StatusCode)
		}
		return LLMResponse{}, fmt.Errorf("API error: %v", errorResponse)
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return LLMResponse{}, fmt.Errorf("error decoding response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return LLMResponse{}, fmt.Errorf("no response generated")
	}

	output := geminiResp.Candidates[0].Content.Parts[0].Text
	if output == "" {
		return LLMResponse{}, fmt.Errorf("empty response received")
	}

	var result LLMResponse
	switch {
	case strings.Contains(output, "##mermaid"):
		result.Diagram = output
	case strings.Contains(output, "table"):
		result.Table = output
	default:
		result.Text = output
	}

	return result, nil
}
