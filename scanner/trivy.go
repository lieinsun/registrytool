package scanner

import (
	"net/http"

	"github.com/aquasecurity/fanal/cache"
	remoteCache "github.com/aquasecurity/trivy/pkg/cache"
	"github.com/aquasecurity/trivy/pkg/rpc/client"
)

type Trivy struct {
	remoteURL client.RemoteURL
	insecure  client.Insecure
	headers   http.Header
	cache     cache.ArtifactCache
	scanner   client.Scanner
}

type Option func(*Trivy)

func New(url string, opts ...Option) *Trivy {
	t := Trivy{
		remoteURL: client.RemoteURL(url),
		headers:   make(http.Header),
	}
	for _, opt := range opts {
		opt(&t)
	}

	remoteScanner := client.NewProtobufClient(t.remoteURL, t.insecure)
	t.scanner = client.NewScanner(client.CustomHeaders(t.headers), remoteScanner)
	t.cache = remoteCache.NewRemoteCache(string(t.remoteURL), t.headers, bool(t.insecure))
	return &t
}
