package auth_reciever

import (
	"authorization_token_repo"
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
	AccessToken  string `json:"access_token"`
	ExpiresIn    int16  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
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

	if resp, err := client.Do(r); err != nil {
		log.Print(err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &accessToken)
		authorization_token_repo.SaveBearer(
			accessToken.AccessToken,
			accessToken.ExpiresIn,
		)

		authorization_token_repo.SaveRefresh(
			accessToken.RefreshToken,
			accessToken.ExpiresIn,
		)
	}

	return accessToken
}
