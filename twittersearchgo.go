// Package twitterquerygo implements a search-optimized Twitter client library in Go.
package twitterquerygo

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

// SearchClient implements a search-optimized Twitter client.
type SearchClient struct {
	Client  twittergo.Client
	SinceID uint64
}

// SearchTweetsResponse implements the response of a search query, containing tweets and the timestamp when the rate limit resets
type SearchTweetsResponse struct {
	Tweets         []twittergo.Tweet
	RateLimitReset time.Time
}

// ISearchClient defines the behaviour of a search-optimized Twitter client.
type ISearchClient interface {
	// SetSinceID Sets the since_id query parameter
	SetSinceID(sinceID uint64)

	// Search searches tweets given a search parameter 'q' till either there are no more results or the rate limit is exceeded
	Search(q string) ([]twittergo.Tweet, error)

	// SearchTillMaxID searches tweets before 'maxID' given a search parameter 'q' till either there are no more results or the rate limit is exceeded
	SearchTillMaxID(q string, maxID uint64) ([]twittergo.Tweet, bool, error)
}

// NewClientUsingApplicationAuth creates a new SearchClient using application authentication, with a rate limited to 450 requests per 15 minutes
func NewClientUsingApplicationAuth(consumerKey string, consumerSecret string) *SearchClient {
	return &SearchClient{
		Client: *twittergo.NewClient(&oauth1a.ClientConfig{
			ConsumerKey:    consumerKey,
			ConsumerSecret: consumerSecret,
		}, nil),
	}
}

// NewClientUsingUserAuth creates a new SearchClient using user authentication, with a rate limited to 180 requests per 15 minutes
func NewClientUsingUserAuth(consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) *SearchClient {
	return &SearchClient{
		Client: *twittergo.NewClient(&oauth1a.ClientConfig{
			ConsumerKey:    consumerKey,
			ConsumerSecret: consumerSecret,
		}, oauth1a.NewAuthorizedConfig(accessToken, accessTokenSecret)),
	}
}

// SetSinceID Sets the since_id query parameter
func (c *SearchClient) SetSinceID(sinceID uint64) {
	c.SinceID = sinceID
}

// Search searches tweets given a search parameter 'q' till either there are no more results or the rate limit is exceeded
func (c *SearchClient) Search(q string) (*SearchTweetsResponse, error) {

	query := url.Values{}
	query.Set("q", q)
	queryURL := fmt.Sprintf("/1.1/search/tweets.json?%v&count=100", query.Encode())
	if c.SinceID > 0 {
		queryURL = fmt.Sprintf("%s&since_id=%d", queryURL, c.SinceID)
	}

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.Client.SendRequest(request)
	if err != nil {
		return nil, err
	}

	results := &twittergo.SearchResults{}
	err = response.Parse(results)
	if err != nil {
		if rateLimitErr, isRateLimitErr := err.(twittergo.RateLimitError); isRateLimitErr {
			return &SearchTweetsResponse{RateLimitReset: rateLimitErr.Reset}, nil
		} else {
			return nil, err
		}
	}

	if results.Statuses() == nil || len(results.Statuses()) == 0 {
		return nil, err
	}

	tweets := results.Statuses()
	var minID uint64 = 18446744073709551615
	for _, tweet := range results.Statuses() {
		if tweet.Id() < minID {
			minID = tweet.Id()
		}
	}

	for {
		moreTweetsResponse, err := c.SearchTillMaxID(q, minID-1)
		if err != nil {
			return nil, err
		}

		if moreTweetsResponse.Tweets == nil || len(moreTweetsResponse.Tweets) == 0 {
			break
		}

		tweets = append(tweets, moreTweetsResponse.Tweets...)
		if time.Now().Before(moreTweetsResponse.RateLimitReset) {
			return &SearchTweetsResponse{
				Tweets:         tweets,
				RateLimitReset: moreTweetsResponse.RateLimitReset,
			}, nil
		}

		for _, tweet := range moreTweetsResponse.Tweets {
			if tweet.Id() < minID {
				minID = tweet.Id()
			}
		}
	}

	return &SearchTweetsResponse{Tweets: tweets}, nil
}

// SearchTillMaxID searches tweets before 'maxID' given a search parameter 'q' till either there are no more results or the rate limit is exceeded
func (c *SearchClient) SearchTillMaxID(q string, maxID uint64) (*SearchTweetsResponse, error) {

	query := url.Values{}
	query.Set("q", q)
	queryURL := fmt.Sprintf("/1.1/search/tweets.json?%v&count=100&max_id=%d", query.Encode(), maxID)
	if c.SinceID > 0 {
		queryURL = fmt.Sprintf("%s&since_id=%d", queryURL, c.SinceID)
	}

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.Client.SendRequest(request)
	if err != nil {
		return nil, err
	}

	results := &twittergo.SearchResults{}
	err = response.Parse(results)
	if err != nil {
		if rateLimitErr, isRateLimitErr := err.(twittergo.RateLimitError); isRateLimitErr {
			return &SearchTweetsResponse{RateLimitReset: rateLimitErr.Reset}, nil
		} else {
			return nil, err
		}
	}

	return &SearchTweetsResponse{
		Tweets: results.Statuses(),
	}, nil
}
