package network

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	endpoint := fmt.Sprintf("%s/%s", c.baseURL, url.PathEscape(key))
	resp, err := c.client.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

// Put sends a PUT request to set a value.
func (c *Client) Put(key, value string) error {
	endpoint := fmt.Sprintf("%s/%s", c.baseURL, url.PathEscape(key))
	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBufferString(value))
	if err != nil {
		return fmt.Errorf("error creating PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making PUT request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	return nil
}

// Delete sends a DELETE request to delete a value.
func (c *Client) Delete(key string) error {
	endpoint := fmt.Sprintf("%s/%s", c.baseURL, url.PathEscape(key))
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
