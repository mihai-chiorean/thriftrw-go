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

// Package pluginapigen provides a plugin Handle that generates code used by the
// plugin system itself.
//
// This is made available as a hidden flag: "--generate-plugin-api".
package pluginapigen

import (
	"path/filepath"
	"strings"

	intplugin "github.com/thriftrw/thriftrw-go/internal/plugin"
	"github.com/thriftrw/thriftrw-go/plugin"
	"github.com/thriftrw/thriftrw-go/plugin/api"
)

// Handle is a plugin.Handle that generates code for the plugin system API.
var Handle intplugin.Handle = handle{}

type handle struct{}

func (handle) Name() string {
	return "pluginapigen"
}

func (handle) Close() error {
	return nil // no-op
}

func (handle) ServiceGenerator() intplugin.ServiceGenerator {
	return sgen{}
}

type sgen struct{}

func (sgen) Handle() intplugin.Handle {
	return Handle
}

func (sgen) Generate(req *api.GenerateServiceRequest) (*api.GenerateServiceResponse, error) {
	files := make(map[string][]byte)
	for _, serviceID := range req.RootServices {
		service := req.Services[serviceID]
		module := req.Modules[service.ModuleID]

		templateData := struct {
			Service *api.Service
			Request *api.GenerateServiceRequest
		}{Service: service, Request: req}

		var (
			err       error
			ifaceOpts []plugin.TemplateOption
		)

		ifacePath := filepath.Join(module.Directory, strings.ToLower(service.Name)+".go")
		ifaceOpts = append(ifaceOpts, templateOptions...)
		ifaceOpts = append(ifaceOpts, plugin.GoFileImportPath(module.Package))
		files[ifacePath], err = plugin.GoFileFromTemplate(
			ifacePath, interfaceTemplate, templateData, ifaceOpts...)
		if err != nil {
			return nil, err
		}

		clientPath := filepath.Join(service.Directory, "client.go")
		files[clientPath], err = plugin.GoFileFromTemplate(
			clientPath, clientTemplate, templateData, templateOptions...)
		if err != nil {
			return nil, err
		}

		handlerPath := filepath.Join(service.Directory, "handler.go")
		files[handlerPath], err = plugin.GoFileFromTemplate(
			handlerPath, handlerTemplate, templateData, templateOptions...)
		if err != nil {
			return nil, err
		}
	}
	return &api.GenerateServiceResponse{Files: files}, nil
}

// convenience function because "index .Request.Services .Service.ParentID"
// doesn't work. Index expects a ServiceID, not *ServiceID.
func getService(req *api.GenerateServiceRequest, id api.ServiceID) *api.Service {
	return req.Services[id]
}

var templateOptions = []plugin.TemplateOption{
	plugin.TemplateFunc("basename", filepath.Base),
	plugin.TemplateFunc("getService", getService),
}

const interfaceTemplate = `
// Code generated by thriftrw --generate-plugin-api
// @generated

<$module := index .Request.Modules .Service.ModuleID>
package <basename $module.Package>

type <.Service.Name> interface {
	<if .Service.ParentID>
		<$parent := getService .Request .Service.ParentID>
		<if eq $parent.ModuleID .Service.ModuleID>
			<$parent.Name>
		<else>
			<$parentModule := index .Request.Modules $parent.ModuleID>
			<import $parentModule.Package>.<$parent.Name>
		<end>
	<end>

	<range .Service.Functions>
		<.Name>(<range .Arguments>
			<.Name> <formatType .Type>,<end>
		) <if .ReturnType>(<formatType .ReturnType>, error)<else>error<end>
	<end>
}
`

const clientTemplate = `
// Code generated by thriftrw --generate-plugin-api
// @generated

package <basename .Service.Package>

<$envelope := import "github.com/thriftrw/thriftrw-go/internal/envelope">
<$module   := import (index .Request.Modules .Service.ModuleID).Package>

// Client implements a <.Service.Name> client.
type client struct {
	<if .Service.ParentID>
		<$parent := getService .Request .Service.ParentID>
		<$parentModule := index .Request.Modules $parent.ModuleID>
		<import $parentModule.Package>.<$parent.Name>
	<end>
	client <$envelope>.Client
}

// NewClient builds a new <.Service.Name> client.
func NewClient(c <$envelope>.Client) <$module>.<.Service.Name> {
	return &client{
		client: c,
		<if .Service.ParentID>
			<$parent := getService .Request .Service.ParentID>
			<$parent.Name>: <import $parent.Package>.NewClient(t),
		<end>
	}
}

<range .Service.Functions>
<$wire := import "github.com/thriftrw/thriftrw-go/wire">

func (c *client) <.Name>(<range .Arguments>
	_<.Name> <formatType .Type>,<end>
) (<if .ReturnType>success <formatType .ReturnType>,<end> err error) {
	args := <.Name>Helper.Args(<range .Arguments>_<.Name>, <end>)

	var body <$wire>.Value
	body, err = args.ToWire()
	if err != nil {
		return
	}

	body, err = c.client.Send("<.ThriftName>", body)
	if err != nil {
		return
	}

	var result <.Name>Result
	if err = result.FromWire(body); err != nil {
		return
	}

	<if .ReturnType>success, <end>err = <.Name>Helper.UnwrapResponse(&result)
	return
}
<end>
`

const handlerTemplate = `
// Code generated by thriftrw --generate-plugin-api
// @generated

package <basename .Service.Package>

<$envelope := import "github.com/thriftrw/thriftrw-go/internal/envelope">
<$module   := import (index .Request.Modules .Service.ModuleID).Package>
<$wire     := import "github.com/thriftrw/thriftrw-go/wire">

// Handler serves an implementation of the <.Service.Name> service.
type Handler struct {
	impl <$module>.<.Service.Name>

	<if .Service.ParentID>
		<$parent := getService .Request .Service.ParentID>
		parent <import $parent.Package>.Handler
	<end>
}

// NewHandler builds a new <.Service.Name> handler.
func NewHandler(service <$module>.<.Service.Name>) Handler {
	return Handler{
		impl: service,
		<if .Service.ParentID>
			<$parent := getService .Request .Service.ParentID>
			parent: <import $parent.Package>.NewHandler(service),
		<end>
	}
}

// Handle receives and handles a request for the <.Service.Name> service.
func (h Handler) Handle(name string, reqValue <$wire>.Value) (<$wire>.Value, error) {
	switch name {
		<range .Service.Functions>
			case "<.ThriftName>":
				var args <.Name>Args
				if err := args.FromWire(reqValue); err != nil {
					return <$wire>.Value{}, err
				}

				result, err := <.Name>Helper.WrapResponse(
					h.impl.<.Name>(<range .Arguments>args.<.Name>, <end>),
				)
				if err != nil {
					return <$wire>.Value{}, err
				}

				return result.ToWire()
		<end>
		default:
			<if .Service.ParentID>
				return h.parent.Handle(name, reqValue)
			<else>
				return <$wire>.Value{}, <$envelope>.ErrUnknownMethod(name)
			<end>
	}
}
`
