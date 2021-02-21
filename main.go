package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/simonknittel/go-oauth-proxy/config"
)

// Redirects the user to auth endpoint
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	state := randomBase64String(16)

	// TODO: Make scope configurable
	w.Header().Add("Location", fmt.Sprintf("%s?client_id=%s&redirect_url=%s&state=%s",
		config.Get("GITHUB_AUTH_ENDPOINT"),
		config.Get("CLIENT_ID"),
		config.Get("REDIRECT_URI"),
		state))

	w.Header().Add("Set-Cookie", fmt.Sprintf("state=%s; HttpOnly=true", state))

	w.WriteHeader(307)
}

// Handles GitHub's callback and redirects user to the frontend with GitHub's
// access token in response body.
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	stateCookie, err := r.Cookie("state")
	if err != nil {
		fmt.Println("Error while reading the state cookie.")
		w.Header().Add("Location", fmt.Sprintf("%s?error=invalid_state_cookie", config.Get("FRONTEND_ENDPOINT")))
		clearStateCookie(w)
		w.WriteHeader(204)
		return
	}

	returnedState := r.Form.Get("state")
	if returnedState != stateCookie.Value {
		fmt.Println("Error while comparing the returned state with the state cookie.")
		w.Header().Add("Location", fmt.Sprintf("%s?error=states_not_matching", config.Get("FRONTEND_ENDPOINT")))
		clearStateCookie(w)
		w.WriteHeader(204)
		return
	}

	code := r.Form.Get("code")
	getAccessToken(code, w)
}

func getAccessToken(code string, w http.ResponseWriter) {
	resp, err := http.PostForm(config.Get("GITHUB_TOKEN_ENDPOINT"), url.Values{
			"client_id": {config.Get("CLIENT_ID")},
			"client_secret": {config.Get("CLIENT_SECRET")},
			"code": {code}})

	if err != nil {
		fmt.Println("Error while while requesting the access token.")
		fmt.Printf("%s\n", err)
		clearStateCookie(w)
		w.WriteHeader(500)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error while reading response body.")
		clearStateCookie(w)
		w.WriteHeader(500)
		return
	}

	redirectToFrontend(body, w)
}

func redirectToFrontend(body []byte, w http.ResponseWriter) {
	base64Body := base64.RawURLEncoding.EncodeToString(body)

	w.Header().Add("Location", fmt.Sprintf("%s?token=%s", config.Get("FRONTEND_ENDPOINT"), base64Body)) // TODO: Not cool. Token should propably not be visible in the URL
	clearStateCookie(w)
	w.WriteHeader(307)
}

func clearStateCookie(w http.ResponseWriter) {
	w.Header().Add("Set-Cookie", "state=; expires=Thu, 01 Jan 1970 00:00:00 GMT")
}

func main() {
	config.Init()
	config.Required()

	http.HandleFunc(config.Get("AUTHORIZE_PATH"), authorizeHandler)
	http.HandleFunc(config.Get("CALLBACK_PATH"), callbackHandler)

	log.Printf("About to listen on port: %s", config.Get("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Get("PORT")), nil))
}
