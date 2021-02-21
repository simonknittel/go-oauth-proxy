package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Redirects the user to auth endpoint
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	state := randomBase64String(16)

	// TODO: Make scope configurable
	w.Header().Add("Location", fmt.Sprintf("%s?client_id=%s&redirect_url=%s&state=%s",
		getEnv("GITHUB_AUTH_ENDPOINT", "https://github.com/login/oauth/authorize"),
		os.Getenv("CLIENT_ID"),
		getEnv("REDIRECT_URI", "http://localhost:8080/callback"),
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
		w.Header().Add("Location", fmt.Sprintf("%s?error=invalid_state_cookie", os.Getenv("FRONTEND_ENDPOINT")))
		clearStateCookie(w)
		w.WriteHeader(204)
		return
	}

	returnedState := r.Form.Get("state")
	if returnedState != stateCookie.Value {
		fmt.Println("Error while comparing the returned state with the state cookie.")
		w.Header().Add("Location", fmt.Sprintf("%s?error=states_not_matching", os.Getenv("FRONTEND_ENDPOINT")))
		clearStateCookie(w)
		w.WriteHeader(204)
		return
	}

	code := r.Form.Get("code")
	getAccessToken(code, w)
}

func getAccessToken(code string, w http.ResponseWriter) {
	resp, err := http.PostForm(getEnv("GITHUB_TOKEN_ENDPOINT", "https://github.com/login/oauth/token"), url.Values{
			"client_id": {os.Getenv("CLIENT_ID")},
			"client_secret": {os.Getenv("CLIENT_SECRET")},
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

	w.Header().Add("Location", fmt.Sprintf("%s?token=%s", os.Getenv("FRONTEND_ENDPOINT"), base64Body)) // TODO: Not cool. Token should propably not be visible in the URL
	clearStateCookie(w)
	w.WriteHeader(307)
}

func clearStateCookie(w http.ResponseWriter) {
	w.Header().Add("Set-Cookie", "state=; expires=Thu, 01 Jan 1970 00:00:00 GMT")
}

// Source: https://stackoverflow.com/a/55860599/3942401
func randomBase64String(l int) string {
	buff := make([]byte, int(math.Round(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}

// Source: https://stackoverflow.com/questions/40326540/how-to-assign-default-value-if-env-var-is-empty
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func main() {
	godotenv.Load()

	if _, ok := os.LookupEnv("CLIENT_ID"); !ok {
		log.Fatal("CLIENT_ID is missing.")
	}

	if _, ok := os.LookupEnv("CLIENT_SECRET"); !ok {
		log.Fatal("CLIENT_SECRET is missing.")
	}

	http.HandleFunc(getEnv("AUTHORIZE_PATH", "/authorize"), authorizeHandler)
	http.HandleFunc(getEnv("CALLBACK_PATH", "/callback"), callbackHandler)

	log.Printf("About to listen on port: %s", getEnv("PORT", "8080"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", getEnv("PORT", "8080")), nil))
}
