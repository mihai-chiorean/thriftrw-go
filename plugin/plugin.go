// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package plugin

import (
	"fmt"
	"log"
	"os"

	"github.com/thriftrw/thriftrw-go/internal/envelope"
	"github.com/thriftrw/thriftrw-go/internal/frame"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/plugin/api/service/plugin"
	"github.com/thriftrw/thriftrw-go/protocol"
)

const _fastPathFrameSize = 10 * 1024 * 1024 // 10 MB

var _proto = protocol.Binary

// Generator provides the Generator feature for ThriftRW>
type Generator interface {
	Generate(*api.GenerateRequest) (*api.GenerateResponse, error)
}

// Plugin defines a plugin.
type Plugin struct {
	Name string

	// If non-nil, Generator generates arbitrary code for services.
	Generator Generator
}

// Main serves the given plugin. It is the entry point to the plugin system.
// User-defined plugins should call Main with their Plugin definition.
func Main(p *Plugin) {
	// The plugin communicates with the ThriftRW process over stdout and stdin
	// of this process. Requests and responses are Thrift envelopes with a
	// 4-byte big-endian encoded length prefix.

	var features []api.Feature
	if p.Generator != nil {
		features = append(features, api.FeatureGenerator)
	}

	server := frame.NewServer(os.Stdin, os.Stdout)
	pluginHandler := plugin.NewHandler(handler{
		server:   server,
		plugin:   p,
		features: features,
	})

	if err := server.Serve(envelope.NewServer(_proto, pluginHandler)); err != nil {
		log.Fatalf("plugin server failed with error: %v", err)
	}
}

type handler struct {
	server   *frame.Server
	plugin   *Plugin
	features []api.Feature
}

func (h handler) Handshake(request *api.HandshakeRequest) (*api.HandshakeResponse, error) {
	return &api.HandshakeResponse{
		Name:       h.plugin.Name,
		ApiVersion: Version,
		Features:   h.features,
	}, nil
}

func (h handler) Goodbye() error {
	h.server.Stop()
	return nil
}

func (h handler) Generate(req *api.GenerateRequest) (*api.GenerateResponse, error) {
	if h.plugin.Generator == nil {
		return nil, fmt.Errorf("I don't implement Generator") // TODO error handling
	}

	return h.plugin.Generator.Generate(req)
}
