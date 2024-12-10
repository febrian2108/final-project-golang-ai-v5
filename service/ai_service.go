// ai_service.go

package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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

    // Log the response status code for debugging
    log.Printf("API response status: %d", resp.StatusCode)

    if resp.StatusCode != http.StatusOK {
        return "", errors.New("failed to analyze data: " + resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
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

    payload := map[string]string{
        "context": context,
        "query":   query,
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

    // Log the response status code for debugging
    log.Printf("API response status: %d", resp.StatusCode)

    if resp.StatusCode != http.StatusOK {
        return model.ChatResponse{}, errors.New("failed to chat with AI: " + resp.Status)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return model.ChatResponse{}, err
    }

    // Parse response JSON
    var responses []model.ChatResponse
    err = json.Unmarshal(body, &responses)
    if err != nil {
        return model.ChatResponse{}, err
    }

    if len(responses) == 0 {
        return model.ChatResponse{}, errors.New("no chat response received")
    }

    return responses[0], nil
}
