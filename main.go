package main

import (
	"net/http"

	"github.com/a7103/github_oauth/config"
)

func main() {
	serv := config.MustGetOauthServ()

	mux := http.NewServeMux()
	mux.HandleFunc("/", serv.HomePage)
	mux.HandleFunc("/oauth/redirect/info", serv.GetUserInfo)
	mux.HandleFunc("/oauth/redirect/emails", serv.GetUserEmails)
	mux.HandleFunc("/oauth/redirect/orgs", serv.GetUserOrgs)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
