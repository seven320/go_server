package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// package twitter

func search_tweet(client *http.Client) ([]Tweet, error) {

	values := url.Values{}
	values.Add("q", "#golang")

	request, err := http.NewRequest("GET", "https://api.twitter.com/1.1/search/tweets.json"+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result Search_response
	json.Unmarshal(b, &result)
	response.Body.Close()

	return result.Statuses, nil
}

func main() {
	config := oauth1.NewConfig(aagaghaohoga, consumerSecret)
	token := oauth1.NewToken(token, tokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	result, err := search_tweet(httpClient)
	if err != nil {
		fmt.Println(err)
	}

	for _, tweet := range result { //デフォルトでは15個のTweetを拾ってくる
		fmt.Println("text:" + tweet.Text) //テキストの表示
	}
}
