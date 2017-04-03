/*
Copyright 2015 ENDOH takanao.
<https://github.com/MiCHiLU/go-lru-cache-stats>

Copyright 2012 Google Inc.
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

// peers.go defines how processes find and communicate with their peers.

package cache

import (
	pb "github.com/golang/groupcache/groupcachepb"
)

// Context is an opaque value passed through calls to the
// ProtoGetter. It may be nil if your ProtoGetter implementation does
// not require a context.
type Context interface{}

// ProtoGetter is the interface that must be implemented by a peer.
type ProtoGetter interface {
	Get(context Context, in *pb.GetRequest, out *pb.GetResponse) error
}
