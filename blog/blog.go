package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	baseBlogURL  = "https://api.hubapi.com/content/api/v2/blog-posts?content_group_id=3708593652&state=published&hapikey="
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

type topicPostData []struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	FeaturedImage string  `json:"featured_image"`
	Slug          string  `json:"slug"`
	PublishDate   int64   `json:"publish_date"`
	TopicIds      []int64 `json:"topic_ids"`
}

// TopicPostsResponseModel . . .
type TopicPostsResponseModel struct {
	Topic topicData     `json:"topic"`
	Posts topicPostData `json:"posts"`
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
	ID   int64  `json:"id"`
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
	TopicIDs      []int64     `json:"topic_ids"`
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

	pageResModel := PageResponseModel{Topics: tdata.Objects, Posts: pdata.Objects}
	if len(data.Objects) > 0 {
		pageResModel.Post = data.Objects[0]
	}
	return &pageResModel, err
}

// GetPostsWithTopicID . . .
func GetPostsWithTopicID(topicSlugString string) (*TopicPostsResponseModel, error) {
	var response TopicPostsResponseModel
	topicsResp, err := http.Get(fmt.Sprintf("%s%s&slug=%s&property=id&property=name&property=slug", baseTopicURL, os.Getenv("hubSpotAPI"), topicSlugString))
	if err != nil {
		return nil, err
	}

	var topicsBody struct {
		Objects []topicData `json:"objects"`
	}
	err = json.NewDecoder(topicsResp.Body).Decode(&topicsBody)
	if err != nil {
		return nil, err
	}

	if topicsBody.Objects == nil || len(topicsBody.Objects) == 0 {
		return nil, fmt.Errorf("could not find topic with slug %s", topicSlugString)
	}
	response.Topic = topicsBody.Objects[0]

	blogPostsResp, err := http.Get(fmt.Sprintf("%s%s&topic_id=%d&property=id&property=name&property=topic_ids&property=featured_image&property=publish_date&property=slug", baseBlogURL, os.Getenv("hubSpotAPI"), response.Topic.ID))
	if err != nil {
		return nil, err
	}

	var blogPostsBody struct {
		Objects topicPostData `json:"objects"`
	}
	if err := json.NewDecoder(blogPostsResp.Body).Decode(&blogPostsBody); err != nil {
		return nil, err
	}
	response.Posts = blogPostsBody.Objects

	return &response, nil
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

	contains := func(t []*topicData, id int64) bool {
		for _, v := range t {
			if v.ID == id {
				return true
			}
		}
		return false
	}

	var t []*topicData
	for _, p := range postsData.Objects {
		for _, postTopicIDs := range p.TopicIDs {
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
		URL:    fmt.Sprintf("%s%s&limit=3&offset=%d&archived=false&property=id&property=html_title&property=post_summary&property=publish_date&property=topic_ids&property=slug&property=featured_image", baseBlogURL, os.Getenv("hubSpotAPI"), offset),
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
		URL:    fmt.Sprintf("%s%s&slug=%s&archived=false&property=featured_image&property=name&property=slug&property=html_title&property=meta_description&property=publish_date&property=post_body&property=blog_author&property=topic_ids", baseBlogURL, os.Getenv("hubSpotAPI"), slug),
		Method: http.MethodGet})
}

func getPosts() ([]byte, error) {
	return doRequest(request{
		URL:    fmt.Sprintf("%s%s&limit=6&archived=false&property=id&property=html_title&property=name&property=post_summary&property=publish_date&property=topic_ids&property=slug&property=featured_image", baseBlogURL, os.Getenv("hubSpotAPI")),
		Method: http.MethodGet})
}

func getFeaturedPosts() ([]byte, error) {
	return doRequest(request{
		URL:    fmt.Sprintf("%s%s&limit=3&archived=false&property=id&property=html_title&property=name&property=slug&property=featured_image", baseBlogURL, os.Getenv("hubSpotAPI")),
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
