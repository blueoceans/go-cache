/*
Copyright 2015 ENDOH takanao.
<https://github.com/MiCHiLU/go-lru-cache-stats>

Copyright 2013 Google Inc.
<https://github.com/golang/groupcache>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lru

import (
	"bytes"
	"net/http"
	"sync"
)

const defaultBasePath = "/_groupcache/"

const defaultReplicas = 50

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// Context optionally specifies a context for the server to use when it
	// receives a request.
	// If nil, the server uses a nil Context.
	Context func(*http.Request) Context

	// Transport optionally specifies an http.RoundTripper for the client
	// to use when it makes a request.
	// If nil, the client uses http.DefaultTransport.
	Transport func(Context) http.RoundTripper

	// base path including leading and trailing slash, e.g. "/_groupcache/"
	basePath string

	// this peer's base URL, e.g. "https://example.net:8000"
	self string

	mu          sync.Mutex             // guards peers and httpGetters
	httpGetters map[string]*httpGetter // keyed by e.g. "http://10.0.0.2:8008"
}

// HTTPPoolOptions are the configurations of a HTTPPool.
type HTTPPoolOptions struct {
	// BasePath specifies the HTTP path that will serve groupcache requests.
	// If blank, it defaults to "/_groupcache/".
	BasePath string

	// Replicas specifies the number of key replicas on the consistent hash.
	// If blank, it defaults to 50.
	Replicas int
}

var httpPoolMade bool

// Set updates the pool's list of peers.
// Each peer value should be a valid base URL,
// for example "http://example.net:8000".
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{transport: p.Transport, baseURL: peer + p.basePath}
	}
}

type httpGetter struct {
	transport func(Context) http.RoundTripper
	baseURL   string
}

var bufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}
