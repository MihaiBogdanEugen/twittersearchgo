// Package twitterquerygo implements a search-optimized Twitter client library in Go.
package twitterquerygo

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

// BatchSize Query for tweets in batches of this size
const BatchSize = 100

// SearchClient implements a search-optimized Twitter client.
type SearchTwitterClient struct {
	TwitterClient twittergo.Client
	SinceID       uint64
	Language      string
}

// SearchTweetsResponse implements the response of a search query, containing tweets and the timestamp when the rate limit resets
type SearchTweetsResponse struct {
	Tweets             []twittergo.Tweet
	HasRateLimit       bool
	RateLimit          uint32
	RateLimitRemaining uint32
	RateLimitReset     time.Time
}

// ISearchClient defines the behaviour of a search-optimized Twitter client.
type ISearchClient interface {
	// SetSinceID sets the since_id query parameter
	SetSinceID(sinceID uint64)

	// SetLang sets the lang query parameter
	SetLanguage(language string)

	// Search searches tweets given a search parameter 'q' till either there are no more results or the rate limit is exceeded
	Search(query string) ([]twittergo.Tweet, error)

	// SearchTillMaxID searches tweets before 'maxID' given a search parameter 'q' till either there are no more results or the rate limit is exceeded
	SearchTillMaxID(query string, maxID uint64) ([]twittergo.Tweet, bool, error)
}

// NewClientUsingApplicationAuth creates a new SearchClient using application authentication, with a rate limited to 450 requests per 15 minutes
func NewClientUsingApplicationAuth(consumerKey string, consumerSecret string) *SearchTwitterClient {
	return &SearchTwitterClient{
		TwitterClient: *twittergo.NewClient(&oauth1a.ClientConfig{
			ConsumerKey:    consumerKey,
			ConsumerSecret: consumerSecret,
		}, nil),
	}
}

// NewClientUsingUserAuth creates a new SearchClient using user authentication, with a rate limited to 180 requests per 15 minutes
func NewClientUsingUserAuth(consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) *SearchTwitterClient {
	return &SearchTwitterClient{
		TwitterClient: *twittergo.NewClient(&oauth1a.ClientConfig{
			ConsumerKey:    consumerKey,
			ConsumerSecret: consumerSecret,
		}, oauth1a.NewAuthorizedConfig(accessToken, accessTokenSecret)),
	}
}

// SetSinceID sets the since_id query parameter
func (c *SearchTwitterClient) SetSinceID(sinceID uint64) {
	c.SinceID = sinceID
}

// SetLang sets the lang query parameter
func (c *SearchTwitterClient) SetLanguage(language string) {
	c.Language = language
}

// Search searches tweets given a search parameter 'q' till either there are no more results or the rate limit is exceeded
func (c *SearchTwitterClient) Search(query string) (*SearchTweetsResponse, error) {

	queryParams := url.Values{}
	queryParams.Set("count", strconv.Itoa(BatchSize))
	if len(c.Language) > 0 {
		queryParams.Set("lang", c.Language)
	}
	queryParams.Set("q", query)
	if c.SinceID > 0 {
		queryParams.Set("since_id", strconv.FormatUint(c.SinceID, 10))
	}
	queryURL := fmt.Sprintf("/1.1/search/tweets.json?%v", queryParams.Encode())

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.TwitterClient.SendRequest(request)
	if err != nil {
		return nil, err
	}

	result := &SearchTweetsResponse{}
	if response.HasRateLimit() {
		result.HasRateLimit = true
		result.RateLimit = response.RateLimit()
		result.RateLimitRemaining = response.RateLimitRemaining()
		result.RateLimitReset = response.RateLimitReset()
	}

	searchResults := &twittergo.SearchResults{}
	if err = response.Parse(searchResults); err != nil {
		if rateLimitErr, isRateLimitErr := err.(twittergo.RateLimitError); isRateLimitErr {
			result.HasRateLimit = true
			result.RateLimit = rateLimitErr.RateLimit()
			result.RateLimitRemaining = rateLimitErr.RateLimitRemaining()
			result.RateLimitReset = rateLimitErr.RateLimitReset()
		} else {
			return nil, err
		}
	}

	if searchResults.Statuses() == nil || len(searchResults.Statuses()) == 0 {
		return result, nil
	}

	result.Tweets = searchResults.Statuses()

	var minID uint64 = 18446744073709551615
	for _, tweet := range searchResults.Statuses() {
		if tweet.Id() < minID {
			minID = tweet.Id()
		}
	}

	for {
		nextResponse, err := c.SearchTillMaxID(query, minID-1)
		if err != nil {
			return nil, err
		}

		result.Tweets = append(result.Tweets, nextResponse.Tweets...)
		result.HasRateLimit = nextResponse.HasRateLimit
		result.RateLimit = nextResponse.RateLimit
		result.RateLimitRemaining = nextResponse.RateLimitRemaining
		result.RateLimitReset = nextResponse.RateLimitReset

		if len(nextResponse.Tweets) == 0 || time.Now().Before(nextResponse.RateLimitReset) {
			break
		}

		for _, tweet := range nextResponse.Tweets {
			if tweet.Id() < minID {
				minID = tweet.Id()
			}
		}
	}

	return result, nil
}

// SearchTillMaxID searches tweets before 'maxID' given a search parameter 'q' till either there are no more results or the rate limit is exceeded
func (c *SearchTwitterClient) SearchTillMaxID(query string, maxID uint64) (*SearchTweetsResponse, error) {

	queryParams := url.Values{}
	queryParams.Set("count", strconv.Itoa(BatchSize))
	queryParams.Set("q", query)
	if len(c.Language) > 0 {
		queryParams.Set("lang", c.Language)
	}
	queryParams.Set("max_id", strconv.FormatUint(maxID, 10))
	if c.SinceID > 0 {
		queryParams.Set("since_id", strconv.FormatUint(c.SinceID, 10))
	}
	queryURL := fmt.Sprintf("/1.1/search/tweets.json?%v", queryParams.Encode())

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.TwitterClient.SendRequest(request)
	if err != nil {
		return nil, err
	}

	result := &SearchTweetsResponse{
		Tweets: []twittergo.Tweet{},
	}

	if response.HasRateLimit() {
		result.HasRateLimit = true
		result.RateLimit = response.RateLimit()
		result.RateLimitRemaining = response.RateLimitRemaining()
		result.RateLimitReset = response.RateLimitReset()
	}

	searchResults := &twittergo.SearchResults{}
	if err = response.Parse(searchResults); err != nil {
		if rateLimitErr, isRateLimitErr := err.(twittergo.RateLimitError); isRateLimitErr {
			result.HasRateLimit = true
			result.RateLimit = rateLimitErr.RateLimit()
			result.RateLimitRemaining = rateLimitErr.RateLimitRemaining()
			result.RateLimitReset = rateLimitErr.RateLimitReset()
		} else {
			return nil, err
		}
	}

	if searchResults.Statuses() != nil && len(searchResults.Statuses()) > 0 {
		result.Tweets = searchResults.Statuses()
	}

	return result, nil
}
