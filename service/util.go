package service

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	ClientIDKey     = "client_id"
	ClientSecretKey = "client_secret"
	RedirectURLKey  = "redirect_uri"
	CodeKey         = "code"
	ScopeKey        = "scope"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func getToken(url string) (Token, error) {
	res, err := Get(nil, url)
	if err != nil {
		return Token{}, err
	}

	var token Token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return Token{}, err
	}

	return token, nil
}

const oauthAccessTokenAPI = "https://github.com/login/oauth/access_token"

func getAccessTokenUrl(config GithubOauthConfig, code string) (string, error) {
	u, err := url.Parse(oauthAccessTokenAPI)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set(ClientIDKey, config.ClientID)
	q.Set(ClientSecretKey, config.ClientSecret)
	q.Set(CodeKey, code)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func getRedirectURL(baseURL string, redirectURLType RedirectURLType) (string, error) {
	u, err := url.JoinPath(baseURL, string(redirectURLType))
	if err != nil {
		return "", err
	}

	return u, nil
}

const oauthAuthorizeURL = "https://github.com/login/oauth/authorize"

func formatOauthAuthorizeURL(config GithubOauthConfig, redirectURLType RedirectURLType, scopes *[]string) (string, error) {
	u, err := url.Parse(oauthAuthorizeURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set(ClientIDKey, config.ClientID)

	redirectURL, err := getRedirectURL(config.RedirectURL, redirectURLType)
	if err != nil {
		return "", err
	}
	q.Set(RedirectURLKey, redirectURL)

	if scopes != nil {
		var stringScopes string
		for _, s := range *scopes {
			stringScopes = fmt.Sprintf("%s %s", stringScopes, s)
		}

		q.Set(ScopeKey, stringScopes)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

type OauthAuthorizeURLs struct {
	UserInfoURL   string
	UserEmailsURL string
	UserOrgsURL   string
}

type RedirectURLType string

const (
	USER_INFO   RedirectURLType = "info"
	USER_EMAILS RedirectURLType = "emails"
	USER_ORGS   RedirectURLType = "orgs"
)

func getOauthAuthorizeURLs(config GithubOauthConfig) (OauthAuthorizeURLs, error) {
	userInfoURL, err := formatOauthAuthorizeURL(config, USER_INFO, nil)
	if err != nil {
		return OauthAuthorizeURLs{}, err
	}

	userEmailsURL, err := formatOauthAuthorizeURL(config, USER_EMAILS, &[]string{"user:email"})
	if err != nil {
		return OauthAuthorizeURLs{}, err
	}

	userOrgsURL, err := formatOauthAuthorizeURL(config, USER_ORGS, &[]string{"read:org"})
	if err != nil {
		return OauthAuthorizeURLs{}, err
	}

	return OauthAuthorizeURLs{
		UserInfoURL:   userInfoURL,
		UserEmailsURL: userEmailsURL,
		UserOrgsURL:   userOrgsURL,
	}, nil
}
