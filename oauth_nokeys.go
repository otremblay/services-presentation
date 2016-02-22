package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var tokenstore map[string]*oauth2.Token = map[string]*oauth2.Token{}

func main() {
	clientId := "YOUR ID" // OMIT
	clientSecret := "YOUR SECRET" // OMIT
	var conf *oauth2.Config
	var githubconf = &oauth2.Config{ // OMIT
		ClientID:     clientId,        // OMIT
		ClientSecret: clientSecret,    // OMIT
		Endpoint:     github.Endpoint, // OMIT
	}
	githubconf = githubconf // OMIT

	var googleconf = &oauth2.Config{
		ClientID:     "YOUR ID", //OMIT
		ClientSecret: "YOUR SECRET",                                                 //OMIT
		RedirectURL:  "http://local.otremblay.com:23232/callback",
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/plus.login"},
	}

	googleconf = googleconf

	conf = googleconf

	http.HandleFunc("/callback", func(rw http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		if code == "" {
			// ...snip...
			rw.WriteHeader(http.StatusBadRequest) // OMIT
			fmt.Fprintln(rw, "No code! Oh noes!") // OMIT
			return                                // OMIT
		}

		tok, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil { // OMIT
			// ...snip...
			rw.WriteHeader(http.StatusTeapot) // OMIT
			fmt.Fprintln(rw, err.Error())     // OMIT
			return                            // OMIT
		}

		cookie := bakeCookie()
		tokenstore[cookie.Value] = tok
		http.SetCookie(rw, cookie)
		http.Redirect(rw, req, "http://local.otremblay.com:23232/", http.StatusTemporaryRedirect)
	})
	// ROOT BEGIN OMIT
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// Thus begins the oauth dance
				oauthDance(conf, rw, req)
				return
			}
			// ...snip...
			rw.WriteHeader(http.StatusTeapot) // OMIT
			fmt.Fprintln(rw, err.Error())     // OMIT
			return                            // OMIT
		}
		tok := tokenstore[c.Value]
		if tok == nil {
			oauthDance(conf, rw, req)
			return
		}
		client := conf.Client(oauth2.NoContext, tok)
		//		resp, err := client.Get("https://api.github.com/user") // OMIT
		resp, err := client.Get("https://www.googleapis.com/plus/v1/people/me")
		if err != nil {
			// ...snip...
			fmt.Fprintln(rw, err) // OMIT
			return                // OMIT
		}
		if resp.StatusCode >= 400 { // OMIT
			b, _ := httputil.DumpResponse(resp, true) // OMIT
			fmt.Println(string(b))                    // OMIT
			return                                    // OMIT
		} // OMIT
		io.Copy(rw, resp.Body)
	})
	//ROOT END OMIT
	http.ListenAndServe(":23232", nil)
}
func bakeCookie() *http.Cookie {
	rb := make([]byte, 32)
	rand.Read(rb)
	rs := base64.URLEncoding.EncodeToString(rb)
	cookie := &http.Cookie{}
	cookie.Name = "token"
	cookie.Value = rs
	return cookie
}

func oauthDance(conf *oauth2.Config, rw http.ResponseWriter, req *http.Request) {
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(rw, req, url, http.StatusTemporaryRedirect)
}
