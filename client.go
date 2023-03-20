package fortytwo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	apiURL     = "https://api.intra.42.fr"
	apiAuth    = "https://api.intra.42.fr/oauth/authorize"
	apiToken   = "https://api.intra.42.fr/oauth/token"
	apiVersion = "v2"
	maxRetries = 3
)

type Token string

func (it Token) String() string {
	return string(it)
}

// ClientOption to configure API client
type ClientOption func(*Client)

type Client struct {
	httpClient      *http.Client
	baseURL         *url.URL
	redirectURL     string
	apiVersion      string
	fortyTwoVersion string
	Scope           []string

	ClientID     string
	ClientSecret string

	Achievement AchievementService
	CursusUser  CursusUserService
	Cursus      CursusService
	User        UserService

	maxRetries int
}

func NewClient(ctx context.Context, ClientID, ClientSecret, RedirectURL string, Scope []string, opts ...ClientOption) (*Client, error) {
	uAPI, err := url.Parse(apiURL)
	if err != nil {
		panic(err)
	}

	cfg := clientcredentials.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes:       []string{"public"},
		TokenURL:     apiToken,
	}

	_, err = cfg.Token(ctx)
	if err != nil {
		return nil, err
	}

	c := &Client{
		redirectURL: RedirectURL,
		httpClient:  cfg.Client(ctx),
		baseURL:     uAPI,
		apiVersion:  apiVersion,
		maxRetries:  maxRetries,
		Scope:       Scope,
	}

	c.Achievement = &AchievementClient{apiClient: c}
	c.CursusUser = &CursusUserClient{apiClient: c}
	c.Cursus = &CursusClient{apiClient: c}
	c.User = &UserClient{apiClient: c}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// WithHTTPClient overrides the default http.Client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithVersion overrides the Intra42 API version
func WithVersion(version string) ClientOption {
	return func(c *Client) {
		c.fortyTwoVersion = version
	}
}

// WithRetry overrides the default number of max retry attempts on 429 errors
func WithRetry(retries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = retries
	}
}

func (c *Client) request(ctx context.Context, method, urlStr, token string, queryParams map[string]string, requestBody interface{}) (*http.Response, error) {
	var err error

	var u *url.URL

	if u, err = c.baseURL.Parse(fmt.Sprintf("%s/%s", c.apiVersion, urlStr)); err != nil {
		return nil, err
	}

	var buf io.ReadWriter

	var body []byte

	if requestBody != nil && !reflect.ValueOf(requestBody).IsNil() {
		body, err = json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(body)
	}

	if len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}

		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Intra42-Version", c.fortyTwoVersion)
	req.Header.Add("Content-Type", "application/json")

	failedAttempts := 0

	var res *http.Response

	for {
		if token != "" {
			res, err = http.DefaultClient.Do(req.WithContext(ctx))
		} else {
			res, err = c.httpClient.Do(req.WithContext(ctx))
		}
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusTooManyRequests {
			break
		}

		failedAttempts++
		if failedAttempts == c.maxRetries {
			return nil, &RateLimitedError{Message: fmt.Sprintf("Retry request with 429 response failed after %d retries", failedAttempts)}
		}
		// https://api.intra.42.fr/apidoc/guides/getting_started#limits
		retryAfterHeader := res.Header["Retry-After"]
		if len(retryAfterHeader) == 0 {
			return nil, &RateLimitedError{Message: "Retry-After header missing from Intra42 API response headers for 429 response"}
		}

		retryAfter := retryAfterHeader[0]

		var waitSeconds int

		if waitSeconds, err = strconv.Atoi(retryAfter); err != nil {
			break // should not happen
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(waitSeconds) * time.Second):
		}
	}

	if res.StatusCode != http.StatusOK {
		var apiErr Error

		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return nil, err
		}
		apiErr.Status = res.StatusCode

		return nil, &apiErr
	}

	return res, nil
}

func (c *Client) GetLink(ctx context.Context, state string) string {
	cfg := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       c.Scope,
		RedirectURL:  c.redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  apiAuth,
			TokenURL: apiToken,
		},
	}

	return cfg.AuthCodeURL(state)
}

func (c *Client) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	cfg := oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       c.Scope,
		RedirectURL:  c.redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  apiAuth,
			TokenURL: apiToken,
		},
	}

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}
