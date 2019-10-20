// package twitter
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("API_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("API_SECRET"))
	return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}

func getUserImage(user string) (imageurl string) {

	loadEnv()
	api := getTwitterApi()
	v := url.Values{}

	User, err := api.GetUsersLookup("yosyuaomenww", v)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", User[0].ProfileImageURL)
	return User[0].ProfileImageURL
}

func main() {
	_ = getUserImage("yosyuaomenww")
}
