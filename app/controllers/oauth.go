package controllers

import (
	"encoding/json"
	"fmt"
	"marketboardproject/app/models"
	"marketboardproject/keys"
	"net/http"

	"github.com/revel/revel"
)

var DiscordUser models.DiscordUser

type Oauth struct {
	*revel.Controller
}

// Will send to Authorization if there is no cookies.
func (c Oauth) Login() revel.Result {
	// Checks if there is a session for this user.
	r := c.Request.In.GetRaw().(*http.Request)
	// If the token expired, we need to renew one.
	usercookie, _ := r.Cookie("access")
	if usercookie != nil {
		// If a session has the token in it still (err == nil), we can safely redirect.
		_, err := c.Session.Get(usercookie.Value)
		if err == nil {
			return c.Redirect("/UserInfo")
		} else {
			return c.Redirect(keys.DiscordURL)
		}
	} else {
		// Otherwise Authorize.
		return c.Redirect(keys.DiscordURL)
	}
}

// Once authorized, discord sends a code in the parameters.
// We add them to the session and put a token to their cookie.
func (c Oauth) User() revel.Result {
	code := c.Params.Get("code")
	accesstokenbytevalue := keys.Oauthparams.DiscordAccessToken(code)
	var access models.AccessToken
	json.Unmarshal(accesstokenbytevalue, &access)
	userbytevalue := keys.Oauthparams.DiscordGetUserObject(access.AccessToken)
	json.Unmarshal(userbytevalue, &DiscordUser)

	// Assign to the session, the discorduser object.
	c.Session[access.AccessToken] = DiscordUser

	// We need to save this access token to a cookie, so that the user can access the information
	accesscookie := &http.Cookie{
		Name:     "access",
		Value:    access.AccessToken,
		Expires:  access.Expiration,
		Secure:   true,
		HttpOnly: true,
	}
	c.SetCookie(accesscookie)
	return c.Redirect("/UserInfo")
}

// Post-Authentication
func (c Oauth) UserInfo() revel.Result {
	// We need to obtain the user's key to get their information.
	r := c.Request.In.GetRaw().(*http.Request)
	usercookie, _ := r.Cookie("access")
	discorduser, _ := c.Session.Get(usercookie.Value)
	// We're gonna need to cast this as a map to get our information back.
	discordmap, _ := discorduser.(map[string]interface{})

	fmt.Println(discordmap)
	t := discordmap["id"]

	return c.Render(t)
}
