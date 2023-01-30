# github oauth
A simple example of the use of github oauth for the Go language

## How to run?
1. Read the [github doc](https://docs.github.com/en/developers/apps/building-oauth-apps/creating-an-oauth-app) for create an Oauth app
2. Set Homepage URL to `http://localhost:8080/`, Set Authorization callback URL to `http://localhost:8080/oauth/redirect` when create Oauth app
   - ![img.png](readmeimg/img.png)
3. Get `Client ID` and `Client Secret`
   - ![client_id and client_secret.png](readmeimg/img_1.png)
4. Populate them into [dev_env](dev_env) or [config.json](config.json)
   - NOTE: 
     - `RedirectURL` = Authorization callback URL (`http://localhost:8080/oauth/redirect`)
     - if `dev_env` and `config.json` simultaneously exist, `dev_env` data will overwrite of `config.json`
5. Open you terminal, execute the following code in order:
      ```go
      go get -v ./...
      ```
 
      ```shell
      # Can be ignored if using config.json
      source dev_env
      ```
 
      ```go
      go run main.go
      ```
6. Now you can do something in `localhost:8080`
   - ![img.png](readmeimg/img_2.png)
