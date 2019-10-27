package twitter

// package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load("./envfiles/.env")
	// err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("API_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("API_SECRET"))
	return anaconda.NewTwitterApi(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
}

func GetUserImage(user string) (imageurl string, err error) {
	loadEnv()
	api := getTwitterApi()
	v := url.Values{}

	User, err := api.GetUsersLookup(user, v)
	if err != nil { //想定できるエラーはapi上限エラーと，存在しないユーザーのエラー
		fmt.Printf("%s", err)
		return "", err
	}
	return User[0].ProfileImageURL, nil
}

// func main() {
// 	userImage, _ := getUserImage("yosyuaomenww")
// 	fmt.Printf(userImage)
// }
