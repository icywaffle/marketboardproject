package controllers

import (
	"encoding/json"
	"marketboardproject/app/models"
	"marketboardproject/keys"

	"github.com/revel/revel"
)

type Oauth struct {
	*revel.Controller
}

// Link to Log in.
func (c Oauth) Login() revel.Result {

	return c.Redirect(keys.DiscordURL)
}

func (c Oauth) User() revel.Result {
	code := c.Params.Get("code")
	accesstokenbytevalue := keys.Oauthparams.DiscordAccessToken(code)
	var access models.AccessToken
	json.Unmarshal(accesstokenbytevalue, &access)
	userbytevalue := keys.Oauthparams.DiscordGetUserObject(access.AccessToken)
	var user models.DiscordUser
	json.Unmarshal(userbytevalue, &user)
	return c.Render(user)
}
