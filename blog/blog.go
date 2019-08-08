package blog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

// GetRecentPosts . . .
func (c *Client) GetRecentPosts() ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&limit=3&archived=false&property=id&property=html_title&property=topic_ids&property=slug&property=featured_image&content_group_id=3708593652&state=published", c.baseURL, c.apiKey),
		Method: http.MethodGet,
	})
}

// GetPost . . .
func (c *Client) GetPost(slug string) ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&slug=%s&archived=false&property=featured_image&property=name&property=slug&property=html_title&property=meta_description&property=publish_date&property=post_body&property=blog_author&state=published&property=topic_ids", c.baseURL, c.apiKey, slug),
		Method: http.MethodGet,
	})
}

// GetTopics . . .
func (c *Client) GetTopics() ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s", c.baseURL, c.apiKey),
		Method: http.MethodGet,
	})
}

// GetPostsWithTopicID . . .
func (c *Client) GetPostsWithTopicID(topicSlugString string) ([]byte, error) {
	var topicSlugID int

	topicSlug := strings.ToLower(topicSlugString)

	ts := strings.ReplaceAll(topicSlug, "-", " ")

	switch ts {
	case "resellers":
	case "dealers":
		topicSlugID = 4463036677
		break
	case "passive das vs active das":
		topicSlugID = 3908520381
		break
	case "small business":
		topicSlugID = 3984404787
		break
	case "office solutions":
		topicSlugID = 3984404787
		break
	case "case studies":
		topicSlugID = 4126584798
		break
	case "m2m":
		topicSlugID = 4463032372
		break
	case "commercial buildings":
		topicSlugID = 4463034582
		break
	case "residential":
		topicSlugID = 4474989641
		break
	case "passive das":
		topicSlugID = 5208231936
		break
	case "event":
	case "events":
		topicSlugID = 5208232018
		break
	case "cellphone signal booster":
	case "cellphone signal boosters":
		topicSlugID = 5232260383
		break
	case "4g signal booster":
		topicSlugID = 5258816479
		break
	case "integrators":
		topicSlugID = 6472497905
		break
	case "financial":
		topicSlugID = 6488598601
		break
	case "healthcare":
		topicSlugID = 6845697866
		break
	case "5g":
		topicSlugID = 8559218284
		break
	default:
		topicSlugID = 0
		break
	}

	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%d?hapikey=%s", c.baseURL, topicSlugID, c.apiKey),
		Method: http.MethodGet,
	})
}

// GetRSS . . .
func (c *Client) GetRSS() ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&limit=5&archived=false&property=id&property=publish_date&property=topic_ids&property=html_title&property=post_summary&property=topic_ids&property=slug&property=featured_image&property=blog_author_id&property=post_body&content_group_id=3708593652&state=published", c.baseURL, c.apiKey),
		Method: http.MethodGet,
	})
}

// GetMorePosts . . .
func (c *Client) GetMorePosts(offset int) ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&limit=3&offset=%d&archived=false&property=id&property=html_title&property=post_summary&property=topic_ids&property=slug&property=featured_image&content_group_id=3708593652&state=published", c.baseURL, c.apiKey, offset),
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
