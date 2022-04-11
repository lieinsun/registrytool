package dockerhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/docker/distribution/reference"
	"github.com/golang/glog"
)

const (
	Domain = "hub.docker.com"
)

// https://www.postman.com/dockerdev/workspace/docker-hub/example/17990590-23905501-eddd-43cd-b624-380188a40835
// /v2/namespaces/{namespace}/repositories/{repository}/images 这个接口可以方便的拿到image及tags的信息 但是需要Pro
// Docker image management features are a Pro & Team feature, to find out more about Docker's Pro and Team subscriptions please go to https://www.docker.com/pricing
const (
	LoginURL            = "/v2/users/login"
	ListRepositoriesURL = "/v2/repositories/%s/"
	ImageDetailURL      = "/v2/repositories/%s/tags/%s/"
	//TagsURL         = "/v2/repositories/%s/tags/"
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
	account string
	image   string
	tag     string
}

func NewClient(opts ...Option) *Client {
	cli := Client{
		client: new(http.Client),
		url: url.URL{
			Scheme: "https",
			Host:   Domain,
		},
	}

	for _, opt := range opts {
		opt(&cli)
	}
	return &cli
}

type Option func(image *Client)

func WithAuth(username, password, token string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
		c.token = token
	}
}

func referencePath(repo, image string) (string, error) {
	s := image
	if repo != "" {
		s = fmt.Sprintf("%s/%s", repo, image)
	}
	ref, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return "", err
	}
	ref = reference.TagNameOnly(ref)
	ref = reference.TrimNamed(ref)
	return reference.Path(ref), nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	glog.Infof("HTTP request: %+v", req)
	req.Header["Accept"] = []string{"application/json"}
	req.Header["Content-Type"] = []string{"application/json"}
	if c.token == "" && req.URL.Path != LoginURL {
		token, err := c.Login(req.Context())
		if err != nil {
			return nil, err
		}
		c.token = token
	}
	req.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", c.token)}

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
		if err == nil {
			if ok, err := extractError(buf, resp); ok {
				return nil, err
			}
		}
		return nil, fmt.Errorf("bad status code %q", resp.Status)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	glog.Infof("HTTP response body: %s", buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// extractError 解析dockerhub特定格式的错误信息
func extractError(buf []byte, resp *http.Response) (bool, error) {
	var responseBody map[string]string
	if err := json.Unmarshal(buf, &responseBody); err == nil {
		for _, k := range []string{"message", "detail"} {
			if msg, ok := responseBody[k]; ok {
				return true, fmt.Errorf("bad status code %q: %s", resp.Status, msg)
			}
		}
	}
	return false, nil
}
