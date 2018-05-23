twittersearchgo
=========
This project extends [kurrik](https://github.com/kurrik)'s [twittergo](https://github.com/kurrik/twittergo) Twitter client with search orientated features.

The goal of this project is to provide an efficient way of working with [timelines](https://developer.twitter.com/en/docs/tweets/timelines/guides/working-with-timelines) by searching for tweets by using max_id and since_id query parameters.
In case of a RateLimitError, the reset of the timestamp will be returned - it's the duty of the caller to wait the required time.

[![Build Status](https://travis-ci.org/MihaiBogdanEugen/twittersearchgo.svg?branch=master)](https://travis-ci.org/MihaiBogdanEugen/twittersearchgo)

Installing
----------
Run

    go get github.com/MihaiBogdanEugen/twittersearchgo

Include in your source:

    import "github.com/MihaiBogdanEugen/twittersearchgo"
    
Using dep:

    dep ensure -add github.com/MihaiBogdanEugen/twittersearchgo

Godoc
-----
See http://godoc.org/github.com/MihaiBogdanEugen/twittersearchgo

Configuration
-----

One can configure a `SearchTwitterClient` by providing at least:
- the Twitter API Consumer Key and
- the Twitter API Consumer Secret.

Thus, the client will use application authentication. For user authentication, the additional following parameters must be provided:
- the Twitter API Access Token and
- the Twitter API Access Token Secret

Apart of these, there are the following optional parameters:
| Name | Required | Description | Default Value | 
| ------------- | ------------- | ------------- | ------------- |
| language | optional | Restricts tweets to the given language, given by an ISO 639-1 code. Language detection is best-effort. | en |
| max_id | optional | Returns results with an ID less than (that is, older than) or equal to the specified ID. | - |
| result_type | optional | Specifies what type of search results you would prefer to receive. Valid values include: `mixed` - Include both popular and real time results in the response; `recent` - return only the most recent results in the response; `popular` - return only the most popular results in the response. | mixed |
| since_id | optional | Returns results with an ID greater than (that is, more recent than) the specified ID. There are limits to the number of Tweets which can be accessed through the API. If the limit of Tweets has occured since the since_id, the since_id will be forced to the oldest ID available. | - |


Credits
-----
All credits go to the original [author](https://github.com/kurrik), this project is a mere extension.

References
-----
- [twittergo](https://github.com/kurrik/twittergo)
- [Standard Search API](https://developer.twitter.com/en/docs/tweets/search/api-reference/get-search-tweets)
- [Get Tweet timelines](https://developer.twitter.com/en/docs/tweets/timelines/guides/working-with-timelines)