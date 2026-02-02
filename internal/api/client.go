package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/OnyxDevTools/onyx-cli/internal/schema"
)

// Client is a minimal HTTP client for the Onyx Schema API.
type Client struct {
	baseURL    string
	databaseID string
	apiKey     string
	apiSecret  string
	http       *http.Client
}

func NewClient(baseURL, databaseID, apiKey, apiSecret string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		databaseID: databaseID,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		http:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) headers() http.Header {
	h := make(http.Header)
	h.Set("x-onyx-key", c.apiKey)
	h.Set("x-onyx-secret", c.apiSecret)
	h.Set("Accept", "application/json")
	h.Set("Content-Type", "application/json")
	return h
}

func (c *Client) request(method, path string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode request: %w", err)
		}
		reader = bytes.NewReader(buf)
	}
	return c.requestRaw(method, path, reader)
}

func (c *Client) requestRaw(method, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header = c.headers()
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s: %s", method, path, strings.TrimSpace(string(data)))
	}
	return data, nil
}

func (c *Client) GetSchema(tables []string) (*schema.SchemaRevision, error) {
	_, rev, err := c.GetSchemaRaw(tables)
	return rev, err
}

func (c *Client) GetSchemaRaw(tables []string) ([]byte, *schema.SchemaRevision, error) {
	var path string
	if len(tables) > 0 {
		trimmed := make([]string, 0, len(tables))
		for _, t := range tables {
			tt := strings.TrimSpace(t)
			if tt != "" {
				trimmed = append(trimmed, url.PathEscape(tt))
			}
		}
		if len(trimmed) > 0 {
			path = fmt.Sprintf("/schemas/%s?tables=%s", url.PathEscape(c.databaseID), strings.Join(trimmed, ","))
		}
	}
	if path == "" {
		path = fmt.Sprintf("/schemas/%s", url.PathEscape(c.databaseID))
	}
	raw, err := c.requestRaw(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	var rev schema.SchemaRevision
	if err := json.Unmarshal(raw, &rev); err != nil {
		return nil, nil, fmt.Errorf("decode schema: %w", err)
	}
	// Normalize entities->tables for callers.
	if len(rev.Tables) == 0 && len(rev.Entities) > 0 {
		rev.Tables = rev.Entities
	}
	if len(rev.Entities) == 0 && len(rev.Tables) > 0 {
		rev.Entities = rev.Tables
	}
	if rev.DatabaseID == "" {
		rev.DatabaseID = c.databaseID
	}
	return raw, &rev, nil
}

func (c *Client) ValidateSchema(req schema.SchemaUpsertRequest) (*schema.SchemaValidationResult, error) {
	body, err := ensureDatabaseIDRaw(req, c.databaseID)
	if err != nil {
		return nil, err
	}
	return c.ValidateSchemaRaw(body)
}

func (c *Client) UpdateSchema(req schema.SchemaUpsertRequest, publish bool) (*schema.SchemaRevision, error) {
	body, err := ensureDatabaseIDRaw(req, c.databaseID)
	if err != nil {
		return nil, err
	}
	return c.UpdateSchemaRaw(body, publish)
}

func (c *Client) ValidateSchemaRaw(body []byte) (*schema.SchemaValidationResult, error) {
	body, err := ensureDatabaseIDRawBytes(body, c.databaseID)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/schemas/%s/validate", url.PathEscape(c.databaseID))
	raw, err := c.requestRaw(http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var res schema.SchemaValidationResult
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, fmt.Errorf("decode validation: %w", err)
	}
	if res.Valid == nil {
		v := true
		res.Valid = &v
	}
	return &res, nil
}

func (c *Client) UpdateSchemaRaw(body []byte, publish bool) (*schema.SchemaRevision, error) {
	body, err := ensureDatabaseIDRawBytes(body, c.databaseID)
	if err != nil {
		return nil, err
	}
	qs := url.Values{}
	if publish {
		qs.Set("publish", "true")
	}
	path := fmt.Sprintf("/schemas/%s", url.PathEscape(c.databaseID))
	if q := qs.Encode(); q != "" {
		path += "?" + q
	}
	raw, err := c.requestRaw(http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	var rev schema.SchemaRevision
	if err := json.Unmarshal(raw, &rev); err != nil {
		return nil, fmt.Errorf("decode schema: %w", err)
	}
	if rev.DatabaseID == "" {
		rev.DatabaseID = c.databaseID
	}
	return &rev, nil
}

// ensureDatabaseIDRaw injects the configured database ID when the payload omits it, while preserving all other fields.
func ensureDatabaseIDRaw(req schema.SchemaUpsertRequest, defaultDB string) ([]byte, error) {
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}
	return ensureDatabaseIDRawBytes(raw, defaultDB)
}

func ensureDatabaseIDRawBytes(raw []byte, defaultDB string) ([]byte, error) {
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("decode schema payload: %w", err)
	}
	if payload["databaseId"] == nil || payload["databaseId"] == "" {
		payload["databaseId"] = defaultDB
	}
	out, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode schema payload: %w", err)
	}
	return out, nil
}
