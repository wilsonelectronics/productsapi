package blog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client . . .
type Client struct {
	baseURL string
	apiKey  string
}

type requestStruct struct {
	URL          string
	Method       string
	OkStatusCode int
}

// NewClient . . .
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// GetAllPosts . . .
func (c *Client) GetAllPosts() ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&limit=1000&property=slug&content_group_id=3708593652", c.baseURL, c.apiKey),
		Method: http.MethodGet,
	})
}

// GetPost . . .
func (c *Client) GetPost(slug string) ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&slug=%s&archived=false&property=id&property=html_title&property=topic_ids&property=slug&property=featured_image&content_group_id=3708593652&state=published", c.baseURL, c.apiKey, slug),
		Method: http.MethodGet,
	})
}

// GetTopics . . .
func (c *Client) GetTopics() ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&property=id&property=name", c.baseURL, c.apiKey),
		Method: http.MethodGet,
	})
}

func (c *Client) doRequest(r requestStruct) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
