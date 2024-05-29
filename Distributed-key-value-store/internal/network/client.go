package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	model "distributed/pkg/models"
)

// Client represents the key-value store client.
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new instance of Client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Get sends a GET request to retrieve a value.
func (c *Client) Get(key string) (string, error) {
	endpoint := fmt.Sprintf("%s/get?key=%s", c.baseURL, url.QueryEscape(key))
	resp, err := c.client.Get(endpoint)
	fmt.Println("getrespose", resp)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var kv model.KeyValue
	if err := json.Unmarshal(body, &kv); err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %w", err)
	}

	return kv.Value, nil
}

// Set sends a POST request to set a value.
func (c *Client) Set(key, value string) error {
	endpoint := fmt.Sprintf("%s/set", c.baseURL)
	kv := model.KeyValue{Key: key, Value: value}
	data, err := json.Marshal(kv)
	if err != nil {
		return fmt.Errorf("error marshaling key-value pair: %w", err)
	}

	resp, err := c.client.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	return nil
}

// Delete sends a DELETE request to delete a value.
func (c *Client) Delete(key string) error {
	endpoint := fmt.Sprintf("%s/delete?key=%s", c.baseURL, url.QueryEscape(key))
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	return nil
}
