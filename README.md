twittersearchgo
=========
This project extends [kurrik](https://github.com/kurrik)'s [twittergo](https://github.com/kurrik/twittergo) Twitter client with search orientated features.

The goal of this project is to provide an efficient way of working with [timelines](https://developer.twitter.com/en/docs/tweets/timelines/guides/working-with-timelines) by searching for tweets by using max_id and since_id query parameters.

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

Credits
-----
All credits go to the original [author](https://github.com/kurrik), this project is a mere extension.

References
-----
- [twittergo](https://github.com/kurrik/twittergo)
- [Standard Search API](https://developer.twitter.com/en/docs/tweets/search/api-reference/get-search-tweets)
- [Get Tweet timelines](https://developer.twitter.com/en/docs/tweets/timelines/guides/working-with-timelines)