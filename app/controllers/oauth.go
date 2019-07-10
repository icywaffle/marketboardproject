package controllers

import (
	"encoding/json"
	"fmt"
	"marketboardproject/app/models"
	"marketboardproject/keys"
	"net/http"
	"time"

	"github.com/revel/revel"
)

var DiscordUser models.DiscordUser

type Oauth struct {
	*revel.Controller
}

// Link to Log in.
func (c Oauth) Login() revel.Result {
	// Checks if there is a session for this user.
	r := c.Request.In.GetRaw().(*http.Request)
	usercookie, _ := r.Cookie("session_token")
	if usercookie != nil {
		if c.Session[usercookie.Value] != nil {
			// If there is a valid cookie, then we just want to go to the info.
			return c.Redirect("/UserInfo")
		} else {
			return c.Redirect(keys.DiscordURL)
		}
	} else {
		// Otherwise Authorize.
		return c.Redirect(keys.DiscordURL)
	}
}

func (c Oauth) User() revel.Result {

	// Do Authorization if the cookie doesn't exist in the session.
	code := c.Params.Get("code")
	accesstokenbytevalue := keys.Oauthparams.DiscordAccessToken(code)
	var access models.AccessToken
	json.Unmarshal(accesstokenbytevalue, &access)
	userbytevalue := keys.Oauthparams.DiscordGetUserObject(access.AccessToken)
	json.Unmarshal(userbytevalue, &DiscordUser)

	c.Session[access.AccessToken] = DiscordUser.UniqueID
	fmt.Println(c.Session)

	// This writes the token to the user cookie.
	w := c.Response.Out.Server.GetRaw().(http.ResponseWriter)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   access.AccessToken,
		Expires: time.Now().Add(120 * time.Second),
	})

	return c.Redirect("/UserInfo")
}

func (c Oauth) UserInfo() revel.Result {
	// When a user access a user only page, then we want to call for it's cookie.
	// Then match it.
	var results interface{}
	r := c.Request.In.GetRaw().(*http.Request)
	w := c.Response.Out.Server.GetRaw().(http.ResponseWriter)
	usercookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusBadRequest)
	}
	sessiontoken := usercookie.Value
	session := c.Session[sessiontoken]
	fmt.Println(session)
	if session != nil {
		results = session
	}
	return c.Render(results)
}
