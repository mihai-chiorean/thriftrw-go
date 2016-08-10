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
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/thriftrw/thriftrw-go/envelope"
	"github.com/thriftrw/thriftrw-go/internal/frame"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/plugin/api/service/plugin"
	"github.com/thriftrw/thriftrw-go/protocol"
	"github.com/thriftrw/thriftrw-go/wire"
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

	Generate(req *api.GenerateRequest) (*api.GenerateResponse, error)
	// TODO(abg): This should be auto generated
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

func (ps multiPlug) Generate(req *api.GenerateRequest) (*api.GenerateResponse, error) {
	// TODO(abg): this should call only those plugins which requested the
	// feature
	var errs []error

	files := make(map[string][]byte)
	usedPaths := make(map[string]string)
	for _, p := range ps {
		res, err := p.Generate(req)
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

	if len(errs) > 0 {
		return nil, multiErr(errs)
	}

	return &api.GenerateResponse{Files: files}, nil
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
	features []api.Feature

	cmd    *exec.Cmd
	stdout io.ReadCloser
	stdin  io.WriteCloser
	client *frame.Client
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

	// TODO(abg): Do we need thread-safety? Plugins will not be called
	// concurrently.

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

	p.client = frame.NewClient(p.stdin, p.stdout)

	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin %q: %v", p.name, err)
	}

	var r plugin.HandshakeResult
	if err := p.send(&plugin.HandshakeArgs{Request: &api.HandshakeRequest{}}, &r); err != nil {
		return fmt.Errorf("handshake failed: %v", err)
	}
	if r.UnsupportedVersionError != nil {
		return r.UnsupportedVersionError
	}
	if r.HandshakeError != nil {
		return r.HandshakeError
	}

	hr := r.Success
	if hr.Name != p.name {
		return fmt.Errorf("plugin name mismatch: expected %q but got %q", p.name, hr.Name)
	}
	if hr.ApiVersion != apiVersion {
		return fmt.Errorf("API version mismatch: expected %q but got %q", apiVersion, hr.ApiVersion)
	}

	p.features = hr.Features // TODO(abg): features should be a map
	p.running = true
	return nil
}

// Generate sends a GenerateRequest to this plugin and gits its response.
func (p *Plugin) Generate(req *api.GenerateRequest) (*api.GenerateResponse, error) {
	if !p.running {
		panic(fmt.Sprintf("Generate(%v) called on plugin %q which is not running", req, p.name))
	}

	var r plugin.GenerateResult
	if err := p.send(&plugin.GenerateArgs{Request: req}, &r); err != nil {
		return nil, err
	}
	if r.GeneratorError != nil {
		return nil, r.GeneratorError
	}
	return r.Success, nil
}

type response interface {
	FromWire(wire.Value) error
}

func (p *Plugin) send(req envelope.Enveloper, res response) error {
	var buff bytes.Buffer
	if err := envelope.Write(_proto, &buff, 1, req); err != nil {
		return err
	}

	resBody, err := p.client.Send(buff.Bytes())
	if err != nil {
		return err
	}

	body, _, err := envelope.ReadReply(_proto, bytes.NewReader(resBody))
	if err != nil {
		return err
	}
	return res.FromWire(body)
}

// Close closes the plugin.
func (p *Plugin) Close() error {
	if !p.running {
		return nil // no-op if already closed
	}

	if p.client != nil {
		if err := p.send(&plugin.GoodbyeArgs{}, &plugin.GoodbyeResult{}); err != nil {
			return err
		}
		p.client = nil
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
