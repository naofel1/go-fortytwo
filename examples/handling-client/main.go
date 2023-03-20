package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/naofel1/go-fortytwo/examples/config"
)

func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	return state
}

func main() {
	ctx := context.Background()

	cfg := &config.API42{
		ClientID:     os.Getenv("FT_API_CLIENT_ID"),
		ClientSecret: os.Getenv("FT_API_CLIENT_SECRET"),
		RedirectURL:  "redirect_url",
		Scopes:       []string{"public"},
	}

	cl := config.Init42API(ctx, cfg)

	// Get a random state (set this as a cookie)
	oauthState := generateStateOauthCookie()

	fmt.Println("Link to authorize application: ", cl.GetLink(ctx, oauthState))

	returnedCode := "test_code"

	// When the user is redirected back from the provider, get the state and the code
	// then check if the state is the same as the one you set in the cookie
	// To get the token from the code, use the following:
	tok, err := cl.GetToken(ctx, returnedCode)
	if err != nil {
		panic(err)
	}

	// Token of client is returned and can be used to make API calls
	fmt.Println("Token: ", tok)
}
