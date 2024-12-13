package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	if table == nil || len(table) == 0 {
		return "", errors.New("table cannot be nil or empty")
	}

	url := "https://api-inference.huggingface.co/models/google/tapas-large-finetuned-wtq"

	payload := map[string]interface{}{
		"table": table,
		"query": query,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	log.Printf("AnalyzeData API response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error response: %s", string(body))
		return "", errors.New("failed to analyze data: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsedResponse struct {
		Cells []string `json:"cells"`
	}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return "", err
	}

	if len(parsedResponse.Cells) == 0 {
		return "", errors.New("no result in API response")
	}

	return parsedResponse.Cells[0], nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	url := "https://api-inference.huggingface.co/models/Qwen/Qwen2.5-Coder-32B-Instruct/v1/chat/completions"

	// API expects a "messages" field with a list of messages
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "system", "content": context}, // System context
			{"role": "user", "content": query},     // User's query
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return model.ChatResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return model.ChatResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, err
	}
	defer resp.Body.Close()

	log.Printf("ChatWithAI API response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error response: %s", string(body))
		return model.ChatResponse{}, errors.New("failed to chat with AI: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ChatResponse{}, err
	}

	// Parse response JSON
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return model.ChatResponse{}, err
	}

	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return model.ChatResponse{}, errors.New("no valid response received from AI")
	}

	return model.ChatResponse{GeneratedText: response.Choices[0].Message.Content}, nil
}
