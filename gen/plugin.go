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

package gen

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/thriftrw/thriftrw-go/internal/envelope"
	"github.com/thriftrw/thriftrw-go/internal/frame"
	"github.com/thriftrw/thriftrw-go/internal/multiplex"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/plugin/api/service/plugin"
	"github.com/thriftrw/thriftrw-go/plugin/api/service/servicegenerator"
	"github.com/thriftrw/thriftrw-go/protocol"
)

const (
	_pluginExecPrefix = "thriftrw-plugin-"
	apiVersion        = "1"
)

var _proto = protocol.Binary

// Plug is the plugin API.
type Plug interface {
	Name() string

	Open() error
	Close() error

	ServiceGenerator() api.ServiceGenerator
}

// Combines a list of Plugs into a single Plug. Requests are sent to all
// plugins and their results are combined.
type multiPlug []Plug

func (ps multiPlug) Open() error {
	var (
		errs   []error
		opened []Plug
	)

	for _, plug := range ps {
		if err := plug.Open(); err != nil {
			errs = append(errs, fmt.Errorf("failed to start plugin %q: %v", plug.Name(), err))
			continue
		}
		opened = append(opened, plug)
	}

	if len(errs) == 0 {
		return nil
	}

	// failed to start all plugins. stop anything that was started.
	for _, p := range opened {
		if err := p.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop plugin %q: %v", p.Name(), err))
		}
	}

	return multiErr(errs)
}

func (ps multiPlug) Close() error {
	var errs []error
	for _, p := range ps {
		if err := p.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop plugin %q: %v", p.Name(), err))
		}
	}
	if len(errs) > 0 {
		return multiErr(errs)
	}
	return nil
}

func (ps multiPlug) ServiceGenerator() api.ServiceGenerator {
	// TODO filter out plugs that don't provide ServiceGenerators
	return multiServiceGenerator(ps)
}

type multiServiceGenerator []Plug

func (ps multiServiceGenerator) Generate(req *api.GenerateServiceRequest) (*api.GenerateServiceResponse, error) {
	var errs []error

	files := make(map[string][]byte)
	usedPaths := make(map[string]string)
	for _, p := range ps {
		sg := p.ServiceGenerator()
		if sg == nil {
			continue
		}
		res, err := sg.Generate(req)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to call plugin %q: %v", p.Name(), err))
			continue
		}

		plugName := p.Name()
		for path, contents := range res.Files {
			usedBy, ok := usedPaths[path]
			if !ok {
				usedPaths[path] = plugName
				files[path] = contents
				continue
			}

			errs = append(errs, fmt.Errorf(
				"plugin conflict: cannot write to %q for plugin %q: "+
					"plugin %q already wrote that file", path, plugName, usedBy))
		}
	}

	// TODO(abg): Validate that none of the GenerateResponses contain ".." paths.
	// If they do, record the name of the plugin that generated that path.

	if len(errs) > 0 {
		return nil, multiErr(errs)
	}

	return &api.GenerateServiceResponse{Files: files}, nil
}

type multiErr []error

func (es multiErr) Error() string {
	msgs := make([]string, len(es)+1)
	msgs[0] = "The following errors occurred:"

	for i, e := range es {
		msgs[i+1] = "  - " + e.Error()
	}

	return strings.Join(msgs, "\n")
}

// Plugin is a code generation plugin for ThriftRW.
type Plugin struct {
	io.Closer

	name string

	running  bool
	path     string
	features map[api.Feature]struct{}

	cmd    *exec.Cmd
	stdout io.ReadCloser
	stdin  io.WriteCloser

	client    api.Plugin
	envClient envelope.Client
}

// NewPlugin builds a new generator plugin.
func NewPlugin(name string) (*Plugin, error) {
	execName := _pluginExecPrefix + name
	path, err := exec.LookPath(execName)
	if err != nil {
		// TODO(abg): This should probably be done at a higher level to validate
		// the list of plugins.
		return nil, fmt.Errorf("invalid plugin %q: executable %q not found on $PATH", name, execName)
	}

	return &Plugin{name: name, path: path}, nil
}

// Name returns the name of the Plugin
func (p *Plugin) Name() string { return p.name }

// Open starts the plugin up.
func (p *Plugin) Open() error {
	if p.running {
		panic(fmt.Sprintf("plugin %q is already running", p.name))
	}

	// TODO(abg): Do we need thread-safety

	// TODO(abg): Maybe plugins will have command line options. Figure out how
	// to get those from users and pass to the plugin.
	p.cmd = exec.Command(p.path)
	p.cmd.Stderr = os.Stderr

	var err error
	p.stdout, err = p.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to allocate stdout pipe: %v", err)
	}

	p.stdin, err = p.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to allocate stdin pipe: %v", err)
	}

	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin %q: %v", p.name, err)
	}

	p.envClient = envelope.NewClient(_proto, frame.NewClient(p.stdin, p.stdout))
	p.client = plugin.NewClient(multiplex.NewClient("Plugin", p.envClient).Send)
	hr, err := p.client.Handshake(&api.HandshakeRequest{})
	if err != nil {
		return fmt.Errorf("handshake failed: %v", err)
	}

	if hr.Name != p.name {
		return fmt.Errorf("plugin name mismatch: expected %q but got %q", p.name, hr.Name)
	}

	if hr.ApiVersion != apiVersion {
		return fmt.Errorf("API version mismatch: expected %q but got %q", apiVersion, hr.ApiVersion)
	}

	p.features = make(map[api.Feature]struct{}, len(hr.Features))
	for _, f := range hr.Features {
		p.features[f] = struct{}{}
	}

	p.running = true
	return nil
}

// ServiceGenerator returns the ServiceGenerator for this plugin or nil if
// this plugin doesn't implement this feature.
func (p *Plugin) ServiceGenerator() api.ServiceGenerator {
	if _, ok := p.features[api.FeatureServiceGenerator]; !ok {
		return nil
	}

	return servicegenerator.NewClient(multiplex.NewClient("ServiceGenerator", p.envClient).Send)
}

// Close closes the plugin.
func (p *Plugin) Close() error {
	if !p.running {
		return nil // no-op if already closed
	}

	if p.client != nil {
		if err := p.client.Goodbye(); err != nil {
			return err
		}
		p.client = nil
		p.envClient = nil
	}

	if p.stdout != nil {
		if err := p.stdout.Close(); err != nil {
			return fmt.Errorf("failed to close stdout stream for plugin %q", p.name)
		}
		p.stdout = nil
	}

	if p.stdin != nil {
		// TODO(abg): Do we need to drain stdin first?
		if err := p.stdin.Close(); err != nil {
			return fmt.Errorf("failed to close stdin stream for plugin %q", p.name)
		}
		p.stdin = nil
	}

	if p.cmd != nil {
		if err := p.cmd.Wait(); err != nil {
			return fmt.Errorf("plugin %q failed: %v", p.name, err)
		}
		p.cmd = nil
	}

	p.running = false
	return nil
}
