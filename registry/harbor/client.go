package harbor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/glog"
)

const (
	CurrentUserURL = "/api/v2.0/users/current"
	ImageDetailURL = "/api/v2.0/projects/%s/repositories/%s/artifacts/%s"
)

type Client struct {
	client   *http.Client
	url      url.URL
	username string
	password string
	token    string
	query
}

type query struct {
	project string
	image   string
	tag     string
}

func NewClient(opts ...Option) *Client {
	cli := Client{
		client: new(http.Client),
		url:    url.URL{},
	}

	for _, opt := range opts {
		opt(&cli)
	}
	return &cli
}

type Option func(image *Client)

func WithHost(host string) Option {
	u, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	return func(c *Client) {
		c.url = *u
	}
}

func WithAuth(username, password, token string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
		c.token = token
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	glog.Infof("HTTP request: %+v", req)
	req.Header["Accept"] = []string{"application/json"}
	req.Header["Content-Type"] = []string{"application/json"}
	if c.token == "" && req.URL.Path != CurrentUserURL {
		c.token, _ = c.Login(req.Context())
	}
	req.Header["Authorization"] = []string{fmt.Sprintf("Basic %s", c.token)}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close() //nolint:errcheck
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, errors.New("unauthorized")
		} else if resp.StatusCode == http.StatusForbidden {
			return nil, errors.New("operation not permitted")
		}
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("bad status code %q: %s", resp.Status, string(buf))
	}
	buf, err := ioutil.ReadAll(resp.Body)
	glog.Infof("HTTP response body: %s", buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
