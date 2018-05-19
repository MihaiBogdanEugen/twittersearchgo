// Package twitterquerygo implements a search-optimized Twitter client library in Go.
package twitterquerygo

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

// SearchClient implements a search-optimized Twitter client.
type SearchClient struct {
	Client  twittergo.Client
	SinceID uint64
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
func (c *SearchClient) Search(q string) ([]twittergo.Tweet, error) {

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
		return nil, err
	}

	var minID uint64 = 18446744073709551615
	var tweets []twittergo.Tweet
	for _, tweet := range results.Statuses() {
		tweets = append(tweets, tweet)
		if tweet.Id() < minID {
			minID = tweet.Id()
		}
	}

	if response.RateLimitRemaining() <= 0 && len(tweets) < 100 {
		return tweets, nil
	}

	for {

		moreTweets, shouldContinue, err := c.SearchTillMaxID(q, minID-1)
		if err != nil {
			return nil, err
		}

		for _, tweet := range moreTweets {
			tweets = append(tweets, tweet)
			if tweet.Id() < minID {
				minID = tweet.Id()
			}
		}

		if !shouldContinue {
			break
		}
	}

	return tweets, nil
}

// SearchTillMaxID searches tweets before 'maxID' given a search parameter 'q' till either there are no more results or the rate limit is exceeded
func (c *SearchClient) SearchTillMaxID(q string, maxID uint64) ([]twittergo.Tweet, bool, error) {

	query := url.Values{}
	query.Set("q", q)
	queryURL := fmt.Sprintf("/1.1/search/tweets.json?%v&count=100&max_id=%d", query.Encode(), maxID)
	if c.SinceID > 0 {
		queryURL = fmt.Sprintf("%s&since_id=%d", queryURL, c.SinceID)
	}

	request, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, false, err
	}

	response, err := c.Client.SendRequest(request)
	if err != nil {
		return nil, false, err
	}

	results := &twittergo.SearchResults{}
	err = response.Parse(results)
	if err != nil {
		return nil, false, err
	}

	var tweets []twittergo.Tweet
	for _, tweet := range results.Statuses() {
		tweets = append(tweets, tweet)
	}

	return tweets, len(tweets) == 100 && response.RateLimitRemaining() > 0, nil
}
