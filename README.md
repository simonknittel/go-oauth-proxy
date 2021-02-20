# go-oauth-proxy

Simple web application written in Go to proxy the OAuth authorization code flow with GitHub in order to prevent a leak of the client secret for frontend applications. This application then could be hosted on something like GCP Functions or similar.

_This project is work in progress and generally more like a experiment to me in order to learn Go._

## Getting started

1. Create a new OAuth app on GitHub: <https://github.com/settings/applications/new>
2. Create a Client secret
3. Duplicate `sample.env` and name it `.env`
4. Add your Client ID and secret to your `.env` file
5. Run the application: `go run .`

## Connect your frontend

1. The the URL of your frontend to the `.env` file
2. On the interface of your frontend, point a link to `http://localhost:8080/authorize`
3. Clicking this link will now redirect your user to GitHub, let them sign in and then redirect them to your frontend. GitHub's response including the access token will be passed as base64 encoded query parameter called `token`.

## Docker

1. Build: `docker build -t simonknittel/go-oauth-proxy:latest .`
2. Run: `docker run --env-file .env -p 8080:8080 simonknittel/go-oauth-proxy:latest`
