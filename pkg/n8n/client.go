package n8n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/renaldis/tutorku-backend/config"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	webhooks   map[string]string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 120 * time.Second},
		baseURL:    config.Cfg.N8NBaseURL,
		webhooks:   config.Cfg.N8NWebhooks,
	}
}

func (c *Client) post(webhook string, payload interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", c.baseURL, c.webhooks[webhook])
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	return result, nil
}

// ── Payload Types ──────────────────────────────────────────

type IngestPayload struct {
	MaterialID string `json:"material_id"`
	UserID     string `json:"user_id"`
	FileBase64 string `json:"file_base64"`
	Filename   string `json:"filename"`
}

type ChatPayload struct {
	MaterialID  string        `json:"material_id"`
	UserID      string        `json:"user_id"`
	Query       string        `json:"query"`
	ChatHistory []ChatHistory `json:"chat_history"`
}

type ChatHistory struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SummarizePayload struct {
	MaterialID string `json:"material_id"`
	UserID     string `json:"user_id"`
	Mode       string `json:"mode"` // short | detailed | mindmap
}

type QuizPayload struct {
	MaterialID string `json:"material_id"`
	UserID     string `json:"user_id"`
	Type       string `json:"type"` // multiple_choice | essay | true_false
	Count      int    `json:"count"`
	Difficulty string `json:"difficulty"` // easy | medium | hard
}

type EssayPayload struct {
	MaterialID string `json:"material_id"`
	UserID     string `json:"user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
}

// ── Methods ────────────────────────────────────────────────

func (c *Client) TriggerIngestion(p IngestPayload) (map[string]interface{}, error) {
	return c.post("ingest", p)
}

func (c *Client) QueryRAG(p ChatPayload) (map[string]interface{}, error) {
	return c.post("chat", p)
}

func (c *Client) Summarize(p SummarizePayload) (map[string]interface{}, error) {
	return c.post("summarize", p)
}

func (c *Client) GenerateQuiz(p QuizPayload) (map[string]interface{}, error) {
	return c.post("quiz", p)
}

func (c *Client) EvaluateEssay(p EssayPayload) (map[string]interface{}, error) {
	return c.post("essay", p)
}
