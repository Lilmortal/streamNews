package twitter

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type bearerToken struct {
	AccessToken string `json:"access_token"`
}

func getBearerToken(client *http.Client) bearerToken {
	consumerKey := os.Getenv("consumerKey")
	consumerSecret := os.Getenv("consumerSecret")

	//Step 1: Encode consumer key and secret
	encodedKeySecret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		url.QueryEscape(consumerKey),
		url.QueryEscape(consumerSecret))))

	//Step 2: Obtain a bearer token
	//The body of the request must be grant_type=client_credentials
	reqBody := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	//The request must be a HTTP POST request
	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", reqBody)
	if err != nil {
		log.Fatal(err, client, req)
	}
	//The request must include an Authorization header formatted as
	//Basic <base64 encoded value from step 1>.
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedKeySecret))
	//The request must include a Content-Type header with
	//the value of application/x-www-form-urlencoded;charset=UTF-8.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("Content-Length", strconv.Itoa(reqBody.Len()))

	//Issue the request and get the bearer token from the JSON you get back
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err, respBody)
	}

	var token bearerToken
	json.Unmarshal(respBody, &token)
	return token
}

// GetTweets does
func GetTweets() {
	client := &http.Client{}

	token := getBearerToken(client)
	//choose your API endpoint that supports application only auth context
	//and create a request object with that
	twitterEndPoint := "https://api.twitter.com/1.1/search/tweets.json?q=donald%20trump"
	req, err := http.NewRequest("GET", twitterEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Step 3: Authenticate API requests with the bearer token
	//include an Authorization header formatted as
	//Bearer <bearer token value from step 2>
	req.Header.Add("Authorization",
		fmt.Sprintf("Bearer %s", token.AccessToken))

	//Issue the request and get the JSON API response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println(string(respBody))
}
