package service

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type GithubOauthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string // Authorization callback URL
}

type OauthServ struct {
	tmpl   *template.Template
	config GithubOauthConfig
}

func MustGetOauthServ(
	tmpl *template.Template,
	config GithubOauthConfig,
) *OauthServ {
	return &OauthServ{
		tmpl:   tmpl,
		config: config,
	}
}

func (o *OauthServ) HomePage(w http.ResponseWriter, r *http.Request) {
	data, err := getOauthAuthorizeURLs(o.config)
	if err != nil {
		panic(err)
	}

	err = o.tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func (o *OauthServ) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	o.oauth(w, r, USER_INFO)
}

func (o *OauthServ) GetUserEmails(w http.ResponseWriter, r *http.Request) {
	o.oauth(w, r, USER_EMAILS)
}

func (o *OauthServ) GetUserOrgs(w http.ResponseWriter, r *http.Request) {
	o.oauth(w, r, USER_ORGS)
}

func (o *OauthServ) oauth(w http.ResponseWriter, r *http.Request, t RedirectURLType) {
	code := r.URL.Query().Get("code")
	url, err := getAccessTokenUrl(o.config, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("get access token url failed, err message: %s", err.Error())))
	}

	token, err := getToken(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("get token failed, err message: %s", err.Error())))

		return
	}

	var res any
	switch t {
	case USER_INFO:
		res, err = getUserInfo(token.AccessToken)
	case USER_EMAILS:
		res, err = getUserEmails(token.AccessToken)
	case USER_ORGS:
		res, err = getUserOrgs(token.AccessToken)
	default:
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(fmt.Sprintf("unsupported method: %s", t)))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("get %s failed, err message: %s", t, err.Error())))
		return
	}

	userInfoBytes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(userInfoBytes)
}

func Get(accessToken *string, URL string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	if accessToken != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accessToken))
	}
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type userInfo struct {
	Login             string    `json:"login"`
	Id                int       `json:"id"`
	AvatarUrl         string    `json:"avatar_url"`
	GravatarId        string    `json:"gravatar_id"`
	Url               string    `json:"url"`
	HtmlUrl           string    `json:"html_url"`
	FollowersUrl      string    `json:"followers_url"`
	FollowingUrl      string    `json:"following_url"`
	GistsUrl          string    `json:"gists_url"`
	StarredUrl        string    `json:"starred_url"`
	SubscriptionsUrl  string    `json:"subscriptions_url"`
	OrganizationsUrl  string    `json:"organizations_url"`
	ReposUrl          string    `json:"repos_url"`
	EventsUrl         string    `json:"events_url"`
	ReceivedEventsUrl string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          string    `json:"hireable"`
	Bio               string    `json:"bio"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

const USER_API = "https://api.github.com/user"

func getUserInfo(accessToken string) (*userInfo, error) {
	res, err := Get(&accessToken, USER_API)
	if err != nil {
		return nil, err
	}

	var userInfo userInfo
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

type userEmail struct {
	Email      string  `json:"email"`
	Primary    bool    `json:"primary"`
	Verified   bool    `json:"verified"`
	Visibility *string `json:"visibility"`
}

const USER_EMAILS_API = "https://api.github.com/user/emails"

func getUserEmails(accessToken string) (*[]userEmail, error) {
	res, err := Get(&accessToken, USER_EMAILS_API)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userEmail []userEmail
	err = json.NewDecoder(res.Body).Decode(&userEmail)
	if err != nil {
		return nil, err
	}

	return &userEmail, nil
}

type userOrgs struct {
	Login            string `json:"login"`
	Id               int    `json:"id"`
	NodeId           string `json:"node_id"`
	Url              string `json:"url"`
	ReposUrl         string `json:"repos_url"`
	EventsUrl        string `json:"events_url"`
	HooksUrl         string `json:"hooks_url"`
	IssuesUrl        string `json:"issues_url"`
	MembersUrl       string `json:"members_url"`
	PublicMembersUrl string `json:"public_members_url"`
	AvatarUrl        string `json:"avatar_url"`
	Description      string `json:"description"`
}

const USER_ORGS_API = "https://api.github.com/user/orgs"

func getUserOrgs(accessToken string) (*[]userOrgs, error) {
	res, err := Get(&accessToken, USER_ORGS_API)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userOrgs []userOrgs
	err = json.NewDecoder(res.Body).Decode(&userOrgs)
	if err != nil {
		return nil, err
	}

	return &userOrgs, nil
}
