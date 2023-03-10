package config

import (
	"github.com/a7103/github_oauth/service"
	"github.com/a7103/github_oauth/views"
	"github.com/jinzhu/configor"
)

var _config *service.GithubOauthConfig

func MustGetGithubOauthConfig() *service.GithubOauthConfig {
	if _config != nil {
		return _config
	}

	cfg := &service.GithubOauthConfig{}
	err := configor.New(&configor.Config{ENVPrefix: "GITHUB_OAUTH", AutoReload: true}).Load(cfg, "config.json")
	if err != nil {
		panic(err)
	}

	_config = cfg

	return _config
}

func MustGetOauthServ() *service.OauthServ {
	return service.MustGetOauthServ(
		views.MustGetTemplate(),
		MustGetGithubOauthConfig(),
	)
}
