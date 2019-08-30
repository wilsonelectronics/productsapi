package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	baseBlogURL  = "https://api.hubapi.com/content/api/v2/blog-posts?hapikey="
	baseTopicURL = "https://api.hubapi.com/blogs/v3/topics?hapikey="
)

type request struct {
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

// CaseStudiesResponseModel . . .
type CaseStudiesResponseModel struct {
	Posts  []*postData  `json:"posts"`
	Topics []*topicData `json:"topic"`
}

// LoadMorePostsResponseModel . . .
type LoadMorePostsResponseModel struct {
	Posts postResponseModel `json:"posts"`
}

// PageResponseModel . . .
type PageResponseModel struct {
	Topics []*topicData `json:"topics"`
	Post   *postData    `json:"post"`
	Posts  []*postData  `json:"featuredPosts"`
}

// HubSpotCookieResponseModel . . .
type HubSpotCookieResponseModel struct {
	VID              int                `json:"vid"`
	CanonicalVid     int                `json:"canonical-vid"`
	MergeVids        []interface{}      `json:"merged-vids"`
	PortalID         int                `json:"portal-id"`
	IsContact        bool               `json:"is-contact"`
	ProfileToken     string             `json:"profile-token"`
	ProfileURL       string             `json:"profile-url"`
	Properties       interface{}        `json:"properties"`
	FormSubmissions  []interface{}      `json:"form-submissions"`
	ListMemberships  []interface{}      `json:"list-memberships"`
	IdentityProfiles []identityProfiles `json:"identity-profiles"`
	MergeAudits      []interface{}      `json:"merge-audits"`
}

type identityProfiles struct {
	VID                     int           `json:"vid"`
	SaveAtTimeStamp         int           `json:"saved-at-timestamp"`
	DeletedChangedTimeStamp int           `json:"deleted-changed-timestamp"`
	Identities              []interface{} `json:"identities"`
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

type postTopic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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

// GetSliderTopicsRecentPosts . . .
func GetSliderTopicsRecentPosts() (*SliderTopicsRecentPosts, error) {
	posts, err := getPosts()
	if err != nil {
		return nil, err
	}

	var pData postResponseModel
	if err = json.Unmarshal(posts, &pData); err != nil {
		return nil, err
	}

	topics, err := getTopics()
	if err != nil {
		return nil, err
	}

	var tData topicResponseModel
	if err = json.Unmarshal(topics, &tData); err != nil {
		return nil, err
	}

	return &SliderTopicsRecentPosts{Posts: pData, Topics: tData}, err
}

// GetPostData . . .
func GetPostData(postSlug string) (*PageResponseModel, error) {
	post, err := getPost(postSlug)
	if err != nil {
		return nil, err
	}

	var data singlePostResponseModel
	if err = json.Unmarshal(post, &data); err != nil {
		return nil, err
	}

	topics, err := getTopics()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(topics))
	var tdata topicResponseModel
	if err = json.Unmarshal(topics, &tdata); err != nil {
		return nil, err
	}

	posts, err := getFeaturedPosts()
	if err != nil {
		return nil, err
	}

	var pdata postResponseModel
	if err = json.Unmarshal(posts, &pdata); err != nil {
		return nil, err
	}

	return &PageResponseModel{Post: data.Objects[0], Topics: tdata.Objects, Posts: pdata.Objects}, err
}

// GetPostsWithTopicID . . .
func GetPostsWithTopicID(topicSlugString string) (*TopicPostsResponseModel, error) {
	var slugID int

	topicSlug := strings.ToLower(topicSlugString)

	ts := strings.ReplaceAll(topicSlug, "-", " ")

	slugID = func() int {
		switch ts {
		case "resellers":
			return 4463036677
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
			return 5208232018
		case "events":
			return 5208232018
		case "cell phone signal booster":
			return 5232260383
		case "cell phone signal boosters":
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

	topic, err := getTopic(slugID)
	if err != nil {
		return nil, err
	}

	var t topicData
	if err = json.Unmarshal(topic, &t); err != nil {
		return nil, err
	}

	posts, err := doRequest(request{
		URL:    fmt.Sprintf("%s%s&limit=1000&property=id&property=name&property=topic_ids&property=featured_image&property=publish_date&property=slug&content_group_id=3708593652&state=published", baseBlogURL, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
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
				tp = append(tp, &postData{FeaturedImage: p.FeaturedImage, HTMLTitle: p.HTMLTitle, ID: p.ID, Name: p.Name, PostSummary: p.PostSummary, Slug: p.Slug, TopicIds: p.TopicIds})
			}
		}
	}

	return &TopicPostsResponseModel{Topic: &t, Posts: tp}, err
}

// GetTwoCaseStudies . . .
func GetTwoCaseStudies() (*CaseStudiesResponseModel, error) {
	posts, err := doRequest(request{
		URL:    fmt.Sprintf("%s%s&limit=2&property=id&property=name&property=topic_ids&property=featured_image&property=publish_date&property=slug&content_group_id=3708593652&state=published&topic_id=4126584798", baseBlogURL, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
	if err != nil {
		return nil, err
	}

	var postsData *postResponseModel
	if err = json.Unmarshal(posts, &postsData); err != nil {
		return nil, err
	}

	topics, err := doRequest(request{
		URL:    fmt.Sprintf("%s%s", baseTopicURL, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
	if err != nil {
		return nil, err
	}

	var tData *topicResponseModel
	if err = json.Unmarshal(topics, &tData); err != nil {
		return nil, err
	}

	contains := func(t []*topicData, id int) bool {
		for _, v := range t {
			if v.ID == id {
				return true
			}
		}
		return false
	}

	var t []*topicData
	for _, p := range postsData.Objects {
		for _, postTopicIDs := range p.TopicIds {
			for _, topicIDs := range tData.Objects {
				if postTopicIDs == topicIDs.ID && !contains(t, topicIDs.ID) {
					t = append(t, &topicData{ID: topicIDs.ID, Name: topicIDs.Name, Slug: topicIDs.Slug})
				}
			}
		}
	}
	return &CaseStudiesResponseModel{Posts: postsData.Objects, Topics: t}, err
}

// LoadMorePosts . . .
func LoadMorePosts(offset int) (*LoadMorePostsResponseModel, error) {
	posts, err := doRequest(request{
		URL:    fmt.Sprintf("%s%s&limit=3&offset=%d&archived=false&property=id&property=html_title&property=post_summary&property=topic_ids&property=slug&property=featured_image&content_group_id=3708593652&state=published", baseBlogURL, os.Getenv("hubSpotAPI"), offset),
		Method: http.MethodGet})
	if err != nil {
		return nil, err
	}

	var pData postResponseModel
	if err = json.Unmarshal(posts, &pData); err != nil {
		return nil, err
	}

	return &LoadMorePostsResponseModel{Posts: pData}, err
}

// GetHubSpotCookies . . .
func GetHubSpotCookies(hubSpotUTK string) (*HubSpotCookieResponseModel, error) {
	cookies, err := doRequest(request{
		URL:    fmt.Sprintf("http://api.hubapi.com/contacts/v1/contact/utk/%s/profile?hapikey=%s&property=form-submissions", hubSpotUTK, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})

	var cookiesData *HubSpotCookieResponseModel
	if err = json.Unmarshal(cookies, &cookiesData); err != nil {
		return nil, err
	}
	return cookiesData, err
}

func getTopic(slugID int) ([]byte, error) {
	return doRequest(request{
		URL:    fmt.Sprintf("https://api.hubapi.com/blogs/v3/topics/%d?hapikey=%s&property=slug", slugID, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
}

func getTopics() ([]byte, error) {
	return doRequest(request{
		URL:    fmt.Sprintf("%s%s", baseTopicURL, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
}

func getPost(slug string) ([]byte, error) {
	return doRequest(request{
		URL:    fmt.Sprintf("%s%s&slug=%s&archived=false&property=featured_image&property=name&property=slug&property=html_title&property=meta_description&property=publish_date&property=post_body&property=blog_author&state=published&property=topic_ids", baseBlogURL, os.Getenv("hubSpotAPI"), slug),
		Method: http.MethodGet})
}

func getPosts() ([]byte, error) {
	return doRequest(request{
		URL:    fmt.Sprintf("%s%s&limit=5&archived=false&property=id&property=html_title&property=post_summary&property=publish_date&property=topic_ids&property=slug&property=featured_image&content_group_id=3708593652&state=published", baseBlogURL, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
}

func getFeaturedPosts() ([]byte, error) {
	return doRequest(request{
		URL:    fmt.Sprintf("%s%s&limit=3&archived=false&property=id&property=html_title&property=slug&property=featured_image&content_group_id=3708593652&state=published", baseBlogURL, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
}

func doRequest(r request) ([]byte, error) {
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
