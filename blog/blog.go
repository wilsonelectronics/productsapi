package blog

import (
	"bytes"
	"encoding/json"
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
	URL    string
	Method string
}

// SliderTopicsRecentPosts . . .
type SliderTopicsRecentPosts struct {
	Topics topicResponseModel `json:"topics"`
	Posts  postResponseModel  `json:"posts"`
}

// TopicPostsResponseModel . . .
type TopicPostsResponseModel struct {
	Topic *topicData  `json:"topic"`
	Posts []*postData `json:"posts"`
}

// PostTopicsResponseModel . . .
type PostTopicsResponseModel struct {
	Topics []*topicData `json:"topics"`
	Post   *postData    `json:"post"`
}

// singlePostResponseModel . . .
type singlePostResponseModel struct {
	Objects []*postData `json:"objects"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
	Total   int         `json:"total"`
}

type blogAuthor struct {
	Avatar            string `json:"avatar"`
	Bio               string `json:"bio"`
	Created           int    `json:"created"`
	DeletedAt         int    `json:"deleted_at"`
	DisplayName       string `json:"display_name"`
	Email             string `json:"email"`
	Facebook          string `json:"facebook"`
	FullName          string `json:"full_name"`
	HasSocialProfiles bool   `json:"has_social_profiles"`
	ID                int    `json:"id"`
	Linkedin          string `json:"linkedin"`
	PortalID          int    `json:"portal_id"`
	Slug              string `json:"slug"`
	Twitter           string `json:"twitter"`
	TwitterUsername   string `json:"twitter_username"`
	Updated           int    `json:"updated"`
	Website           string `json:"website"`
}

type topicData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type postResponseModel struct {
	Limit      int         `json:"limit"`
	Objects    []*postData `json:"objects"`
	Offset     int         `json:"offset"`
	Total      int         `json:"total"`
	TotalCount int         `json:"totalCount,omitempty"`
}

type topicResponseModel struct {
	Objects []*topicData `json:"objects"`
}

type postData struct {
	FeaturedImage string      `json:"featured_image"`
	HTMLTitle     string      `json:"html_title"`
	Name          string      `json:"name"`
	ID            int         `json:"id"`
	PostSummary   string      `json:"post_summary"`
	PostBody      string      `json:"post_body"`
	PublishDate   int         `json:"publish_date"`
	Slug          string      `json:"slug"`
	TopicIds      []int       `json:"topic_ids"`
	BlogAuthor    *blogAuthor `json:"blog_author"`
}

// NewClient . . .
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// GetSliderTopicsRecentPosts . . .
func (c *Client) GetSliderTopicsRecentPosts() (*SliderTopicsRecentPosts, error) {

	posts, err := c.getPosts()
	if err != nil {
		return nil, err
	}

	var pData postResponseModel
	if err = json.Unmarshal(posts, &pData); err != nil {
		return nil, err
	}

	topics, err := c.getTopics()
	if err != nil {
		return nil, err
	}

	var tData topicResponseModel
	if err = json.Unmarshal(topics, &tData); err != nil {
		return nil, err
	}

	sliderTopicRecentPostsData := &SliderTopicsRecentPosts{
		Posts:  pData,
		Topics: tData,
	}

	return sliderTopicRecentPostsData, err
}

// GetPostData . . .
func (c *Client) GetPostData(postSlug string) (*PostTopicsResponseModel, error) {
	post, err := c.getPost(postSlug)
	if err != nil {
		return nil, err
	}

	var data singlePostResponseModel
	if err = json.Unmarshal(post, &data); err != nil {
		return nil, err
	}

	fmt.Println(data)

	singlePost := &PostTopicsResponseModel{
		Topics: nil,
		Post:   data.Objects[0],
	}
	return singlePost, err
}

func (c *Client) getTopic(slugID int) ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("https://api.hubapi.com/blogs/v3/topics/%d?hapikey=%s&property=slug", slugID, c.apiKey),
		Method: http.MethodGet,
	})
}

func (c *Client) getTopics() ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("https://api.hubapi.com/blogs/v3/topics?hapikey=%s&property=id&property=name&property=slug", c.apiKey),
		Method: http.MethodGet,
	})
}

func (c *Client) getPost(slug string) ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&slug=%s&archived=false&property=featured_image&property=name&property=slug&property=html_title&property=meta_description&property=publish_date&property=post_body&property=blog_author&state=published&property=topic_ids", c.baseURL, c.apiKey, slug),
		Method: http.MethodGet,
	})
}

// GetPostsWithTopicID . . .
func (c *Client) GetPostsWithTopicID(topicSlugString string) (*TopicPostsResponseModel, error) {
	var slugID int

	topicSlug := strings.ToLower(topicSlugString)

	ts := strings.ReplaceAll(topicSlug, "-", " ")

	slugID = func() int {
		switch ts {
		case "resellers":
		case "dealers":
			return 4463036677
		case "passive das vs active das":
			return 3908520381
		case "small business":
			return 3984404787
		case "office solutions":
			return 3984404787
		case "case studies":
			return 4126584798
		case "m2m":
			return 4463032372
		case "commercial buildings":
			return 4463034582
		case "residential":
			return 4474989641
		case "passive das":
			return 5208231936
		case "event":
		case "events":
			return 5208232018
		case "cellphone signal booster":
		case "cellphone signal boosters":
			return 5232260383
		case "4g signal booster":
			return 5258816479
		case "integrators":
			return 6472497905
		case "financial":
			return 6488598601
		case "healthcare":
			return 6845697866
		case "5g":
			return 8559218284
		}
		return 0
	}()

	topic, err := c.getTopic(slugID)
	if err != nil {
		return nil, err
	}

	var t topicData
	if err = json.Unmarshal(topic, &t); err != nil {
		return nil, err
	}

	posts, err := c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&limit=1000&property=id&property=name&property=topic_ids&property=featured_image&property=publish_date&property=slug&content_group_id=3708593652&state=published", c.baseURL, c.apiKey),
		Method: http.MethodGet,
	})
	if err != nil {
		return nil, err
	}
	var tPosts *postResponseModel

	if err = json.Unmarshal(posts, &tPosts); err != nil {
		return nil, err
	}

	var tp []*postData
	for _, p := range tPosts.Objects {
		for _, id := range p.TopicIds {
			if id == slugID {
				newPost := &postData{FeaturedImage: p.FeaturedImage, HTMLTitle: p.HTMLTitle, ID: p.ID, Name: p.Name, PostSummary: p.PostSummary, Slug: p.Slug, TopicIds: p.TopicIds}
				tp = append(tp, newPost)
			}
		}
	}

	topicPosts := &TopicPostsResponseModel{
		Topic: &t,
		Posts: tp,
	}

	return topicPosts, err
}

func (c *Client) getPosts() ([]byte, error) {
	return c.doRequest(requestStruct{
		URL:    fmt.Sprintf("%s%s&limit=5&archived=false&property=id&property=html_title&property=post_summary&property=publish_date&property=topic_ids&property=slug&property=featured_image&content_group_id=3708593652&state=published", c.baseURL, c.apiKey),
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
