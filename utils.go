package main

import (
	"crypto/rand"
	"encoding/base64"
	"math"
	"net/http"
)

// Source: https://stackoverflow.com/a/55860599/3942401
func randomBase64String(l int) string {
	buff := make([]byte, int(math.Round(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}

func clearStateCookie(w http.ResponseWriter) {
	w.Header().Add("Set-Cookie", "state=; expires=Thu, 01 Jan 1970 00:00:00 GMT")
}
