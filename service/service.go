package service

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type GithubOauthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string // Authorization callback URL
}

type OauthServ struct {
	tmpl   *template.Template
	config *GithubOauthConfig
}

func MustGetOauthServ(
	tmpl *template.Template,
	config *GithubOauthConfig,
) *OauthServ {
	return &OauthServ{
		tmpl:   tmpl,
		config: config,
	}
}

func (o *OauthServ) HomePage(w http.ResponseWriter, r *http.Request) {
	data, err := getOauthAuthorizeURLs(*o.config)
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

const (
	USER_API        = "https://api.github.com/user"
	USER_EMAILS_API = "https://api.github.com/user/emails"
	USER_ORGS_API   = "https://api.github.com/user/orgs"
)

func (o *OauthServ) oauth(w http.ResponseWriter, r *http.Request, t RedirectURLType) {
	code := r.URL.Query().Get("code")
	url, err := getAccessTokenUrl(*o.config, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("get access token url failed, err message: %s", err.Error())))
		return
	}

	token, err := getToken(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("get token failed, err message: %s", err.Error())))
		return
	}

	var resp *http.Response
	switch t {
	case USER_INFO:
		resp, err = Get(&token.AccessToken, USER_API)
	case USER_EMAILS:
		resp, err = Get(&token.AccessToken, USER_EMAILS_API)
	case USER_ORGS:
		resp, err = Get(&token.AccessToken, USER_ORGS_API)
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
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
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
