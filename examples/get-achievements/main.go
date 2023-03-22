package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/naofel1/go-fortytwo"
	"github.com/naofel1/go-fortytwo/examples/config"
)

func main() {
	ctx := context.Background()

	cfg := &config.API42{
		ClientID:     os.Getenv("FT_API_CLIENT_ID"),
		ClientSecret: os.Getenv("FT_API_CLIENT_SECRET"),
		RedirectURL:  "redirect_url",
		Scopes:       []string{"public"},
	}

	cl := config.Init42API(ctx, cfg)

	achievements, _, err := cl.Achievement.List(ctx, &fortytwo.AchievementQueryRequest{
		Pagination: &fortytwo.Pagination{
			Cursor:   1,
			PageSize: 10,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(achievements)
	if err != nil {
		log.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err := json.Indent(buf, b, "", "\t"); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
