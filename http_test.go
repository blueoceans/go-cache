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
	"flag"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	peerAddrs = flag.String("test_peer_addrs", "", "Comma-separated list of peer addresses; used by TestHTTPPool")
	peerIndex = flag.Int("test_peer_index", -1, "Index of which peer this child is; used by TestHTTPPool")
	peerChild = flag.Bool("test_peer_child", false, "True if running as a child process; used by TestHTTPPool")
)

func testKeys(n int) (keys []string) {
	keys = make([]string, n)
	for i := range keys {
		keys[i] = strconv.Itoa(i)
	}
	return
}

// This is racy. Another process could swoop in and steal the port between the
// call to this function and the next listen call. Should be okay though.
// The proper way would be to pass the l.File() as ExtraFiles to the child
// process, and then close your copy once the child starts.
func pickFreeAddr(t *testing.T) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}

func addrToURL(addr []string) []string {
	url := make([]string, len(addr))
	for i := range addr {
		url[i] = "http://" + addr[i]
	}
	return url
}

func awaitAddrReady(t *testing.T, addr string, wg *sync.WaitGroup) {
	defer wg.Done()
	const max = 1 * time.Second
	tries := 0
	for {
		tries++
		c, err := net.Dial("tcp", addr)
		if err == nil {
			c.Close()
			return
		}
		delay := time.Duration(tries) * 25 * time.Millisecond
		if delay > max {
			delay = max
		}
		time.Sleep(delay)
	}
}
