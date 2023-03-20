package fortytwo

import (
	"io"
	"log"
)

func closeBody(body io.ReadCloser) {
	if errClose := body.Close(); errClose != nil {
		log.Println("failed to close body, should never happen")
	}
}
