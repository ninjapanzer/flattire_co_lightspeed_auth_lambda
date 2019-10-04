package auth_reciever

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type AuthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int64 `json:"expires_in"`
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type Message struct {
	Name string
	Body string
	Time int64
}

func ExchangeCode(code string) (accessToken AccessTokenResponse) {

	apiUrl := "https://cloud.lightspeedapp.com"
	resource := "/oauth/access_token.php"
	data := url.Values{}
	data.Set("client_id", os.Getenv("lightspeed_client_id"))
	data.Set("client_secret", os.Getenv("lightspeed_client_secret"))
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String() // "https://api.com/user/"

	log.Print(urlStr)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)

	//resp, err := http.Post("https://cloud.lightspeedapp.com/oauth/access_token.php", "application/x-www-form-urlencoded", bytes.NewBuffer(payload))
	log.Print(err, resp)
	body, err := ioutil.ReadAll(resp.Body)
	log.Print(resp.Body, err)
	mes := json.Unmarshal([]byte(body), &accessToken)

	log.Print(mes)
	log.Print(body)
	return
}